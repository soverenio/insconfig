package insconfig

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/soverenio/vanilla/throw"

	"github.com/soverenio/insconfig/utils"
)

// Empty config generation from struct part

// you may use insconfig:"comment" Tag on struct fields to express your feelings.
// default values goes from provided obj

type YamlTemplatableStruct interface {
	TemplateTo(w io.Writer, m YamlTemplaterStruct) error
}

type YamlTemplaterStruct struct {
	Obj   interface{}       // what are we marshaling right now
	Level int               // Level of recursion
	Tag   reflect.StructTag // Tag for current field
	FName string            // current field name

	componentPath     utils.Path // path to processable field
	prefixWithSpace   bool
	prefixWithNewline bool
}

func NewYamlTemplaterStruct(obj interface{}) YamlTemplaterStruct {
	return YamlTemplaterStruct{
		Obj:             obj,
		Level:           -1,
		Tag:             "",
		componentPath:   ".",
		prefixWithSpace: false,
	}
}

func getIndent(level int) string {
	if level > 0 {
		return strings.Repeat("  ", level)
	}
	return ""
}

func writeComment(w io.Writer, indent string, tag reflect.StructTag) error {
	icTag, ok := tag.Lookup("insconfig")
	if !ok {
		return nil
	}

	var (
		splitTag = strings.SplitN(icTag, "|", 2)
		comment  = strings.TrimSpace(splitTag[len(splitTag)-1])
	)

	if comment == "" {
		return nil
	}

	if _, err := fmt.Fprintf(w, "%s# %s\n", indent, comment); err != nil {
		return throw.W(err, "failed to write field comment")
	}
	return nil
}

func writeName(w io.Writer, indent string, name string, parentTag reflect.StructTag, objectType reflect.Type) error {
	nameOverride, ok := parentTag.Lookup("yaml")
	if ok {
		name = nameOverride
	}

	if name == "" || objectType.Kind() == reflect.Array {
		return nil
	}

	if _, err := fmt.Fprintf(w, "%s%s: ", indent, strings.ToLower(name)); err != nil {
		return throw.W(err, "failed to write field name")
	}
	return nil
}

func (m YamlTemplaterStruct) TemplateTo(w io.Writer) (err error) {
	if _, ok := w.(*utils.Trimmer); !ok {
		trimmer := utils.NewTrimmer(w)
		defer func() {
			if err != nil {
				return
			}

			err = trimmer.Close()
		}()
		w = trimmer
	}

	if o, ok := m.Obj.(YamlTemplatableStruct); ok {
		return o.TemplateTo(w, m)
	}

	var (
		t      = reflect.TypeOf(m.Obj)
		v      = reflect.ValueOf(m.Obj)
		indent = getIndent(m.Level)
	)

	if t.Kind() == reflect.Ptr {
		m.Obj = v.Elem().Interface()
		return m.TemplateTo(w)
	}

	if err := writeComment(w, indent, m.Tag); err != nil {
		return err
	}

	if err := writeName(w, indent, m.FName, m.Tag, t); err != nil {
		return err
	}

	switch t.Kind() {
	case reflect.Struct:
		if m.Level > 0 || m.prefixWithNewline { // initial struct, no need for newline here
			if _, err := fmt.Fprint(w, "\n"); err != nil {
				return err
			}
		}
		for i := 0; i < t.NumField(); i++ {
			fldType := t.Field(i)
			childTemplater := YamlTemplaterStruct{
				Obj:   v.Field(i).Interface(),
				Level: m.Level + 1,
				Tag:   fldType.Tag,
				FName: fldType.Name,

				componentPath:     m.componentPath.AppendStructKey(fldType.Name),
				prefixWithNewline: true,
			}

			err := childTemplater.TemplateTo(w)
			if err != nil {
				return throw.W(err, "failed to write child element", struct {
					FieldName utils.Path
				}{FieldName: childTemplater.componentPath})
			}
		}

	case reflect.Map:
		if m.prefixWithSpace { // initial struct, no need for newline here
			if _, err := fmt.Fprint(w, " "); err != nil {
				return err
			}
		}

		if _, err := fmt.Fprintf(w, "# <map> of %s\n", t.Elem()); err != nil {
			return err
		}

		iter := v.MapRange()
		for iter.Next() {
			mapKey := fmt.Sprintf("%s", iter.Key().Interface())

			if _, err := fmt.Fprintf(w, "%s  %s:", indent, mapKey); err != nil {
				return err
			}

			childTemplater := YamlTemplaterStruct{
				Obj:             iter.Value().Interface(),
				Level:           m.Level + 1,
				componentPath:   m.componentPath.AppendMapKey(mapKey),
				prefixWithSpace: true,
			}

			err := childTemplater.TemplateTo(w)
			if err != nil {
				return throw.W(err, "failed to write child element", struct {
					FieldName utils.Path
				}{FieldName: childTemplater.componentPath})
			}
		}

	case reflect.Array, reflect.Slice:
		if m.prefixWithSpace { // initial struct, no need for newline here
			if _, err := fmt.Fprint(w, " "); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(w, "# <array> of %s\n", t.Elem()); err != nil {
			return err
		}

		for i := 0; i < v.Len(); i++ {
			childTemplater := YamlTemplaterStruct{
				Obj:   v.Index(i).Interface(),
				Level: m.Level + 1,

				componentPath:   m.componentPath.AppendArrayIdx(i),
				prefixWithSpace: true,
				// prefixWithNewline: true,
			}

			if _, err := fmt.Fprintf(w, "%s  -", indent); err != nil {
				return err
			}
			if err := childTemplater.TemplateTo(w); err != nil {
				return throw.W(err, "failed to write child element", struct {
					FieldName utils.Path
				}{FieldName: childTemplater.componentPath})
			}
		}

	case reflect.String, reflect.Bool, reflect.Uintptr,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:

		if m.prefixWithSpace { // initial struct, no need for newline here
			if _, err := fmt.Fprint(w, " "); err != nil {
				return err
			}
		}

		if _, err := fmt.Fprintf(w, "%v # %s\n", v.Interface(), t.Name()); err != nil {
			return throw.W(err, "failed to write scalar", struct {
				FieldName utils.Path
			}{FieldName: m.componentPath})
		}

	default:
		return throw.New("unknown serialization (please implement YamlTemplatableStruct)", struct {
			Type string
			Kind string
		}{
			Type: t.Name(),
			Kind: t.Kind().String(),
		})
	}
	return nil
}
