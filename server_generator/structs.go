// Package main ...
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

func findStructPairList(path string, endpointParams []string) (map[string]*PackageStructPair, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return nil, err
	}

	structList := findStructList(pkgs)
	structPair := make(map[string]*PackageStructPair, 0)
	for _, s := range structList {
		structEndpointParams := make([]string, len(endpointParams))
		copy(structEndpointParams, endpointParams)

		filePath := fset.File(s.StructObject.Struct).Name()
		fileName := filepath.Base(filePath)
		fileName = fileName[:len(fileName)-len(filepath.Ext(filePath))]

		controllerName := s.StructName
		var structMode StructMode

		if strings.HasSuffix(controllerName, "Request") {
			controllerName = controllerName[:len(controllerName)-7]
			structMode = StructModeRequest
		} else if strings.HasSuffix(controllerName, "Response") {
			controllerName = controllerName[:len(controllerName)-8]
			structMode = StructModeResponse
		} else {
			continue
		}

		if _, ok := structPair[controllerName]; !ok {
			structPair[controllerName] = new(PackageStructPair)
		}
		structPair[controllerName].FileName = filePath

		if strings.HasPrefix(fileName, "0_") {
			fileName = fileName[2:]
			if !strings.HasPrefix(strings.ToLower(fileName), ":id") {
				fileName = strings.ReplaceAll(fileName, "_id", "ID")
			}
			endpoint := strcase.ToCamel(fileName)
			param := strings.ToLower(string(endpoint[0])) + endpoint[1:]
			structEndpointParams = append(structEndpointParams, param)
		}

		switch structMode {
		case StructModeRequest:
			structPair[controllerName].Request = s
			if len(structEndpointParams) > 0 {
				structPair[controllerName].LastParam = structEndpointParams[len(structEndpointParams)-1]
			}

			var requestParams []RequestParam
			requestParams, err = createRequestParams(s.StructName, fset, s.StructObject, structEndpointParams)
			if err != nil {
				return nil, err
			}
			structPair[controllerName].Request.RequestParams = requestParams

			if err = validateRequestByEndpointParams(fset, s.StructObject, structEndpointParams); err != nil {
				return nil, err
			}

			if strings.HasPrefix(strings.ToLower(controllerName), "get") {
				if err = validateGetRequestTags(fset, s.StructObject, structEndpointParams); err != nil {
					return nil, err
				}
			}
		case StructModeResponse:
			structPair[controllerName].Response = s
		}
	}

	return structPair, err
}

func validateRequestByEndpointParams(fset *token.FileSet, structType *ast.StructType, endpointParams []string) error {
	fieldList := structType.Fields.List
	hasEndpoints := make(map[string]bool)
	for _, e := range endpointParams {
		if _, ok := hasEndpoints[e]; !ok {
			hasEndpoints[e] = false
		}
	}

	for _, fields := range fieldList {
		if len(fields.Names) > 1 {
			return fmt.Errorf("%+v: 同じ行に複数のパラメータを記述することはできません。", fields.Names)
		}

		var fieldName string
		if len(fields.Names) > 0 {
			fieldName = fields.Names[0].Name
		} else if t, ok := fields.Type.(*ast.SelectorExpr); ok {
			fieldName = t.Sel.Name
		}

		if fields.Tag != nil {
			tags := reflect.StructTag(strings.Trim(fields.Tag.Value, "`"))
			if paramTag, ok := tags.Lookup("param"); ok {
				fieldName = strings.Split(paramTag, ",")[0]
			}
		}

		if _, ok := hasEndpoints[fieldName]; ok {
			hasEndpoints[fieldName] = true
		}
	}

	var requireParams string
	for name, e := range hasEndpoints {
		if e {
			continue
		}
		if len(requireParams) > 0 {
			requireParams += ", "
		}
		requireParams += name
	}

	if len(requireParams) == 0 {
		return nil
	}

	return fmt.Errorf("%s: Pathマッチング用のパラメータが不足しています。不足しているパラメータ[%s]",
		fset.Position(structType.Pos()).String(), requireParams)
}

func validateGetRequestTags(fset *token.FileSet, structType *ast.StructType, endpointParams []string) error {
	fieldList := structType.Fields.List
	ep := make(map[string]struct{})
	for _, e := range endpointParams {
		if _, ok := ep[e]; !ok {
			ep[e] = struct{}{}
		}
	}

	for i := range fieldList {
		if fieldList[i].Tag == nil {
			continue
		}

		tags := reflect.StructTag(strings.Trim(fieldList[i].Tag.Value, "`"))

		paramTag, ok := tags.Lookup("param")
		if ok {
			_, ok = ep[paramTag]
			if ok {
				continue
			}
		}

		jsonTag, ok := tags.Lookup("json")

		if !ok {
			continue
		}

		queryTag, ok := tags.Lookup("query")

		jsonName := strings.TrimSpace(strings.Split(jsonTag, ",")[0])
		queryName := strings.TrimSpace(strings.Split(queryTag, ",")[0])

		if !ok || jsonName != queryName {
			return fmt.Errorf(
				`%s (json:"%s"!=query:"%s"): GETのRequest structはjsonタグをqueryタグに同じ値を指定する必要があります`,
				fset.Position(fieldList[i].Tag.Pos()).String(),
				jsonName, queryName,
			)
		}
	}

	return nil
}

func createRequestParams(structName string, fset *token.FileSet, st *ast.StructType, endpointParams []string) ([]RequestParam, error) {
	ep := make(map[string]struct{})
	for _, e := range endpointParams {
		if _, ok := ep[e]; !ok {
			ep[e] = struct{}{}
		}
	}

	fl := st.Fields.List
	result := make([]RequestParam, len(fl))

	for i, f := range fl {
		if len(f.Names) > 1 {
			return nil, fmt.Errorf("%s: %+v: 同じ行に複数のパラメータを記述することはできません。",
				fset.Position(f.Pos()).String(), f.Names)
		}

		var fDataType string
		switch types.ExprString(f.Type) {
		case "string":
			fDataType = StringRequestDataType
		case "int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64":
			fDataType = IntegerRequestDataType
		case "float32", "float64":
			fDataType = NumberRequestDataType
		case "bool":
			fDataType = BooleanRequestDataType
		default:
			fDataType = types.ExprString(f.Type)
		}

		fType, fName := getFieldNameFromStructAndEndpointParams(structName, f, ep)

		isRequired := fType == PathRequestName
		if !isRequired && f.Tag != nil {
			tags := reflect.StructTag(strings.Trim(f.Tag.Value, "`"))
			if validateTags, ok := tags.Lookup("validate"); ok {
				isRequired = strings.Contains(validateTags, "required")
			}
		}

		result[i] = RequestParam{
			Name:     fName,
			Type:     fType,
			DataType: fDataType,
			Comment:  parseComment(f.Comment),
			Required: isRequired,
		}
	}

	return result, nil
}

func getFieldNameFromStructAndEndpointParams(structName string, f *ast.Field, ep map[string]struct{}) (RequestParamType, string) {
	var (
		targetTag string
		fName     string
		tags      reflect.StructTag
		fType     RequestParamType
	)

	if len(f.Names) > 0 {
		fName = f.Names[0].Name
	} else if t, ok := f.Type.(*ast.SelectorExpr); ok {
		fName = t.Sel.Name
	}

	if f.Tag != nil {
		tags = reflect.StructTag(strings.Trim(f.Tag.Value, "`"))
		if nameFromTag, ok := tags.Lookup("param"); ok {
			nameFromTag = strings.Split(nameFromTag, ",")[0]
			if _, ok := ep[nameFromTag]; ok {
				return PathRequestName, nameFromTag
			}
		}
	}

	if _, ok := ep[fName]; ok {
		return PathRequestName, fName
	}

	if strings.HasPrefix(strings.ToLower(structName), "get") {
		fType = QueryRequestName
		targetTag = "query"
	} else {
		fType = BodyRequestName
		targetTag = "json"
	}

	if f.Tag == nil {
		return fType, fName
	}

	if nameFromTag, ok := tags.Lookup(targetTag); ok {
		fName = nameFromTag
	}

	return fType, fName
}

var epMap = make(map[string]string)

func findStructList(pkgs map[string]*ast.Package) []*PackageStruct {
	structList := make([]*PackageStruct, 0)

	for _, pkg := range pkgs {
		for fileName, f := range pkg.Files {
			objects := f.Scope.Objects
			for _, object := range objects {
				if object.Kind != ast.Typ {
					continue
				}

				name := object.Name
				if strings.HasSuffix(name, "Request") {
					epMap[name] = ""
				}

				tSpec, ok := object.Decl.(*ast.TypeSpec)
				if !ok {
					continue
				}

				t := tSpec.Type
				switch structObject := t.(type) {
				case *ast.StructType:
					structList = append(structList, &PackageStruct{
						FileName:     fileName,
						PackageName:  pkg.Name,
						StructName:   name,
						StructObject: structObject,
					})
				case *ast.Ident:
					//log.Printf("<IDENT> %s (%s)", name, t.(*ast.Ident).Name)
				default:
					//log.Printf("name=%s , %#v\n", name, tSpec)
				}
			}
		}
	}

	return structList
}

func parseComment(comments *ast.CommentGroup) (comment string) {
	if comments == nil {
		return comment
	} else if len(comments.List) == 0 {
		return comment
	}

	for _, c := range comments.List {
		comment += strings.TrimPrefix(c.Text, "// ")
	}

	return comment
}
