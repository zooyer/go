package properties

import (
	"io"
	"io/ioutil"
	"strings"
	"errors"
	"encoding/xml"
)

const Header = `<!DOCTYPE properties SYSTEM "http://github.com/PublicGiter">` + "\n"

type XmlSupport interface {
	load(props *Properties, in io.Reader) error
	store(props *Properties, out io.Writer, comment, encoding string) error
}

var PROVIDER = newXmlSupport()

func load(props *Properties, in io.Reader) error {
	return PROVIDER.load(props, in)
}

func save(props *Properties, out io.Writer, comment, encoding string) error {
	return PROVIDER.store(props, out, comment, encoding)
}


type xmlSupport struct { }

func (x *xmlSupport) load(props *Properties, in io.Reader) error {
	data,err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	var m XMLProperties
	if err = xml.Unmarshal(data, &m); err != nil {
		return err
	}

	p := props.New()
	if err = x.toProperties(p, &m); err != nil {
		return err
	}

	props.Hashtable = p.Hashtable

	return nil
}

func (x *xmlSupport) store(props *Properties, out io.Writer, comment, encoding string) error {
	switch strings.ToUpper(encoding) {
	case "UTF-8":
	default:
		return errors.New("not support encoding <" + encoding + ">")
	}

	xp,err := x.toXML(props)
	if err != nil {
		return err
	}

	data,err := xp.ToXML(comment)
	if err != nil {
		return err
	}

	_,err = out.Write(data)

	return err
}

func (x *xmlSupport) toProperties(props *Properties, xp *XMLProperties) error {
	if xp == nil {
		return errors.New("XMLProperties is <nil>")
	}

	for _,v := range xp.Entry {
		props.Put(v.Key, v.CDATA)
	}

	return nil
}

func (x *xmlSupport) toXML(props *Properties) (*XMLProperties, error) {
	var xp = new(XMLProperties)
	xp.Entry = make([]XMLReaderEntry, props.Size())
	for i,key := range props.StringPropertyNames() {
		xp.Entry[i] = XMLReaderEntry{
			Key   : key,
			CDATA : props.GetProperty(key),
		}
	}

	return xp, nil
}



func newXmlSupport() XmlSupport {
	return new(xmlSupport)
}

type XMLReaderEntry struct {
	XMLName     xml.Name    `xml:"entry"`
	Key         string      `xml:"key,attr"`
	CDATA       string      `xml:",cdata"`
}

type XMLWriterEntry struct {
	XMLName     xml.Name    `xml:"entry"`
	Key         string      `xml:"key,attr"`
	Data        string      `xml:",chardata"`
}

type XMLReader struct {
	XMLName     xml.Name              `xml:"properties"`
	Comment     string                `xml:"comment"`
	Entry       []XMLReaderEntry      `xml:"entry"`
}

type XMLWriter struct {
	XMLName     xml.Name              `xml:"properties"`
	Entry       []XMLWriterEntry      `xml:"entry"`
}

type XMLWriterComment struct {
	Comment     string                `xml:"comment"`
	XMLWriter
}

type XMLProperties struct {
	XMLReader
}

func (p *XMLProperties) ToXML(comments string) ([]byte, error) {
	w := new(XMLWriter)
	w.XMLName = p.XMLName
	w.Entry = make([]XMLWriterEntry, len(p.Entry))
	for i,e := range p.Entry {
		w.Entry[i] = XMLWriterEntry{
			XMLName  : e.XMLName,
			Key      : e.Key,
			Data     : e.CDATA,
		}
	}

	if comments == "" {
		return xml.MarshalIndent(w, "", "\t")
	}

	wc := new(XMLWriterComment)
	wc.XMLWriter = *w
	wc.Comment = comments

	return xml.MarshalIndent(wc, "", "\t")
}