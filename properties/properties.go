package properties

import (
	"sync"
	"io"
	"bytes"
	"bufio"
	"runtime"
	"time"
	"errors"
	"reflect"
)

// A table of hex digits
var hexDigit  = []byte{ '0','1','2','3','4','5','6','7','8','9','A','B','C','D','E','F' }

// The {@code Properties} class represents a persistent set of
// properties. The {@code Properties} can be saved to a stream
// or loaded from a stream. Each key and its corresponding value in
// the property list is a string.
// <p>
// A property list can contain another property list as its
// "defaults"; this second property list is searched if
// the property key is not found in the original property list.
// <p>
// Because {@code Properties} inherits from {@code Hashtable}, the
// {@code put} and {@code putAll} methods can be applied to a
// {@code Properties} object.  Their use is strongly discouraged as they
// allow the caller to insert entries whose keys or values are not
// {@code Strings}.  The {@code setProperty} method should be used
// instead.  If the {@code store} or {@code save} method is called
// on a "compromised" {@code Properties} object that contains a
// non-{@code String} key or value, the call will fail. Similarly,
// the call to the {@code propertyNames} or {@code list} method
// will fail if it is called on a "compromised" {@code Properties}
// object that contains a non-{@code String} key.
//
// <p>
// The {@link #load(java.io.Reader) load(Reader)} <tt>/</tt>
// {@link #store(java.io.Writer, java.lang.String) store(Writer, String)}
// methods load and store properties from and to a character based stream
// in a simple line-oriented format specified below.
//
// The {@link #load(java.io.InputStream) load(InputStream)} <tt>/</tt>
// {@link #store(java.io.OutputStream, java.lang.String) store(OutputStream, String)}
// methods work the same way as the load(Reader)/store(Writer, String) pair, except
// the input/output stream is encoded in ISO 8859-1 character encoding.
// Characters that cannot be directly represented in this encoding can be written using
// Unicode escapes as defined in section 3.3 of
// <cite>The Java&trade; Language Specification</cite>;
// only a single 'u' character is allowed in an escape
// sequence. The native2ascii tool can be used to convert property files to and
// from other character encodings.
//
// <p> The {@link #loadFromXML(InputStream)} and {@link
// #storeToXML(OutputStream, String, String)} methods load and store properties
// in a simple XML format.  By default the UTF-8 character encoding is used,
// however a specific encoding may be specified if required. Implementations
// are required to support UTF-8 and UTF-16 and may support other encodings.
// An XML properties document has the following DOCTYPE declaration:
//
// <pre>
// &lt;!DOCTYPE properties SYSTEM "http://java.sun.com/dtd/properties.dtd"&gt;
// </pre>
// Note that the system URI (http://java.sun.com/dtd/properties.dtd) is
// <i>not</i> accessed when exporting or importing properties; it merely
// serves as a string to uniquely identify the DTD, which is:
// <pre>
//    &lt;?xml version="1.0" encoding="UTF-8"?&gt;
//
//    &lt;!-- DTD for properties --&gt;
//
//    &lt;!ELEMENT properties ( comment?, entry* ) &gt;
//
//    &lt;!ATTLIST properties version CDATA #FIXED "1.0"&gt;
//
//    &lt;!ELEMENT comment (#PCDATA) &gt;
//
//    &lt;!ELEMENT entry (#PCDATA) &gt;
//
//    &lt;!ATTLIST entry key CDATA #REQUIRED&gt;
// </pre>
//
// <p>This class is thread-safe: multiple threads can share a single
// <tt>Properties</tt> object without the need for external synchronization.
//
// @see <a href="../../../technotes/tools/solaris/native2ascii.html">native2ascii tool for Solaris</a>
// @see <a href="../../../technotes/tools/windows/native2ascii.html">native2ascii tool for Windows</a>
//
// @author  Zhongyuan Zhang
// @since   Golang1.9.2
type Properties struct {
	Hashtable
	mutex     sync.Mutex

	// A property list that contains default values for any keys not
	// found in this property list.
	defaults  *Properties
}

// Creates an empty property list with no default values.
func NewProperties() *Properties {
	var hash = NewHashtable()
	hash.Init()
	return &Properties{
		Hashtable : hash,
		defaults  : nil,
	}
}

// Creates an empty property list with the specified defaults.
func NewPropertiesDefault(defaults *Properties) *Properties {
	var properties = NewProperties()
	properties.defaults = defaults

	return properties
}

// Calls the <tt>Hashtable</tt> method {@code put}. Provided for
// parallelism with the <tt>getProperty</tt> method. Enforces use of
// strings for property keys and values. The value returned is the
// result of the <tt>Hashtable</tt> call to {@code put}.
func (p *Properties) SetProperty(key, value string) interface{} {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.Put(key, value)
}

// The specified stream remains open after this method returns.
// @param   reader   the input character stream.
// @since   1.6
func (p *Properties) Load(reader io.Reader) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.load0(NewLineReader(reader))
}

func (p *Properties) load0(lr LineReader) error {
	var convtBuf = make([]byte, 4096)
	var limit, keyLen, valueStart int
	var c byte
	var hasSep, precedingBackslash bool

	for limit = lr.readLine(); limit >= 0; limit = lr.readLine() {
		c = 0
		keyLen = 0
		valueStart = limit
		hasSep = false

		//fmt.Println("line=<" + string(lr.(*lineReader).lineBuf[:limit]) + ">")
		precedingBackslash = false
		for keyLen < limit {
			c = lr.(*lineReader).lineBuf[keyLen]
			//need check if escaped.
			if (c == '=' || c == ':') && !precedingBackslash {
				valueStart = keyLen + 1
				hasSep = true
				break
			} else if (c == ' ' || c == '\t' ||  c == '\f') && !precedingBackslash {
				valueStart = keyLen + 1
				break
			}
			if c == '\\' {
				precedingBackslash = !precedingBackslash
			} else {
				precedingBackslash = false
			}
			keyLen++
		}
		for valueStart < limit {
			c = lr.(*lineReader).lineBuf[valueStart]
			if c != ' ' && c != '\t' &&  c != '\f' {
				if !hasSep && (c == '=' ||  c == ':') {
					hasSep = true
				} else {
					break
				}
			}
			valueStart++
		}

		key,err   := p.loadConvert(lr.(*lineReader).lineBuf, 0, keyLen, convtBuf)
		if err != nil {
			return err
		}
		value,err := p.loadConvert(lr.(*lineReader).lineBuf, valueStart, limit - valueStart, convtBuf)
		if err != nil {
			return err
		}
		p.Put(key, value)
	}

	return nil
}

// Converts encoded &#92;uxxxx to unicode chars
// and changes special saved chars to their original forms
func (p *Properties) loadConvert(in []byte, off, length int, convtBuf []byte) (string, error) {
	if len(convtBuf) < length {
		var newLen = length * 2
		if newLen < 0 {
			newLen = int(^uint(0) >> 1)
		}
		convtBuf = make([]byte, newLen)
	}

	var aChar byte
	var out = convtBuf
	var outLen = 0
	var end = off + length

	for off < end {
		aChar = in[off]
		off++
		if aChar == '\\' {
			aChar = in[off]
			off++
			if aChar == 'u' {
				// Read the xxxx
				var value = 0

				for i := 0; i < 4; i++ {
					aChar = in[off]
					off++
					switch aChar {
					case '0','1','2','3','4','5','6','7','8','9':
					value = (value << 4) + int(aChar) - '0'
					case 'a','b','c','d','e','f':
					value = (value << 4) + 10 + int(aChar) - 'a'
					case 'A','B','C','D','E','F':
					value = (value << 4) + 10 + int(aChar) - 'A'
					default:
						return "", errors.New("malformed \\uxxxx encoding")
					}
				}
				out[outLen] = byte(value)
				outLen++
			} else {
				if aChar == 't'{ aChar = '\t' } else
				if aChar == 'r'{ aChar = '\r'} else
				if aChar == 'n'{ aChar = '\n'} else
				if aChar == 'f'{ aChar = '\f'}
				out[outLen] = aChar
				outLen++
			}
		} else {
			out[outLen] = aChar
			outLen++
		}
	}

	return string(out[:outLen]), nil
}

// Converts unicodes to encoded &#92;uxxxx and escapes
// special characters with a preceding slash
func (p *Properties) saveConvert(theString string, escapeSpace, escapeUnicode bool) string {
	var length = len(theString)
	var bufLen = length * 2
	if bufLen < 0 {
		bufLen = int(^uint(0) >> 1)
	}

	var outBuffer bytes.Buffer

	for x := 0; x < length; x++ {
		var aChar = theString[x]
		// Handle common case first, selecting largest block that
		// avoids the specials below
		if (aChar > 61) && (aChar < 127) {
			if aChar == '\\' {
				outBuffer.WriteByte('\\')
				outBuffer.WriteByte('\\')
				continue
			}
			outBuffer.WriteByte(aChar)
			continue
		}
		switch aChar {
		case ' ':
			if x == 0 || escapeSpace {
				outBuffer.WriteByte('\\')
			}
			outBuffer.WriteByte(' ')
		case '\t':
			outBuffer.WriteByte('\\')
			outBuffer.WriteByte('t')
		case '\n':
			outBuffer.WriteByte('\\')
			outBuffer.WriteByte('n')
		case '\r':
			outBuffer.WriteByte('\\')
			outBuffer.WriteByte('r')
		case '\f':
			outBuffer.WriteByte('\\')
			outBuffer.WriteByte('f')
		case '=':
			fallthrough
		case ':','#','!':
			outBuffer.WriteByte('\\')
			outBuffer.WriteByte(aChar)
		default:
			if ((aChar < 0x0020) || (aChar > 0x007e)) && escapeUnicode {
				outBuffer.WriteByte('\\')
				outBuffer.WriteByte('u')
				outBuffer.WriteByte(toHex(int(aChar >> 12) & 0xF))
				outBuffer.WriteByte(toHex(int(aChar >>  8) & 0xF))
				outBuffer.WriteByte(toHex(int(aChar >>  4) & 0xF))
				outBuffer.WriteByte(toHex(int(aChar)       & 0xF))
			} else {
				outBuffer.WriteByte(aChar)
			}
		}
	}

	return outBuffer.String()
}

// Calls the {@code store(OutputStream out, String comments)} method
// and suppresses IOExceptions that were thrown.
//
// @deprecated This method does not throw an IOException if an I/O error
// occurs while saving the property list.  The preferred way to save a
// properties list is via the {@code store(OutputStream out,
// String comments)} method or the
// {@code storeToXML(OutputStream os, String comment)} method.
//
// @param   out      an output stream.
// @param   comments   a description of the property list.
// @exception  ClassCastException  if this {@code Properties} object
//             contains any keys or values that are not
//             {@code Strings}.
// Deprecated
func (p *Properties) Save(writer io.Writer, comments string) {
	p.Store(writer, comments)
}

// Writes this property list (key and element pairs) in this
// {@code Properties} table to the output character stream in a
// format suitable for using the {@link #load(java.io.Reader) load(Reader)}
// method.
// <p>
// Properties from the defaults table of this {@code Properties}
// table (if any) are <i>not</i> written out by this method.
// <p>
// If the comments argument is not null, then an ASCII {@code #}
// character, the comments string, and a line separator are first written
// to the output stream. Thus, the {@code comments} can serve as an
// identifying comment. Any one of a line feed ('\n'), a carriage
// return ('\r'), or a carriage return followed immediately by a line feed
// in comments is replaced by a line separator generated by the {@code Writer}
// and if the next character in comments is not character {@code #} or
// character {@code !} then an ASCII {@code #} is written out
// after that line separator.
// <p>
// Next, a comment line is always written, consisting of an ASCII
// {@code #} character, the current date and time (as if produced
// by the {@code toString} method of {@code Date} for the
// current time), and a line separator as generated by the {@code Writer}.
// <p>
// Then every entry in this {@code Properties} table is
// written out, one per line. For each entry the key string is
// written, then an ASCII {@code =}, then the associated
// element string. For the key, all space characters are
// written with a preceding {@code \} character.  For the
// element, leading space characters, but not embedded or trailing
// space characters, are written with a preceding {@code \}
// character. The key and element characters {@code #},
// {@code !}, {@code =}, and {@code :} are written
// with a preceding backslash to ensure that they are properly loaded.
// <p>
// After the entries have been written, the output stream is flushed.
// The output stream remains open after this method returns.
// <p>
//
// @param   writer      an output character stream writer.
// @param   comments   a description of the property list.
// @exception  IOException if writing this property list to the specified
//             output stream throws an <tt>IOException</tt>.
// @exception  ClassCastException  if this {@code Properties} object
//             contains any keys or values that are not {@code Strings}.
// @exception  NullPointerException  if {@code writer} is null.
// @since 1.6
func (p *Properties) Store(writer io.Writer, comments string) error {
	return p.store0(bufio.NewWriter(writer), comments, true)
}

func (p *Properties) store0(bw *bufio.Writer, comments string, escUnicode bool) (err error) {
	if comments != "" {
		if err = writeComments(bw, comments); err != nil {
			return err
		}
	}
	if _,err = bw.WriteString("# " + time.Now().Format(time.UnixDate)); err != nil {
		return err
	}
	if _,err = bw.Write(newLine()); err != nil {
		return err
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _,key := range p.Keys() {
		var val = p.Get(key)

		var skey = key.(string)
		var sval = val.(string)

		skey = p.saveConvert(skey, true, escUnicode)
		// No need to escape embedded and trailing spaces for value, hence
		// pass false to flag.
		sval = p.saveConvert(sval, false, escUnicode)
		if _,err = bw.WriteString(skey + " = " + sval); err != nil {
			return err
		}
		if _,err = bw.Write(newLine()); err != nil {
			return err
		}
	}

	return bw.Flush()
}

// Loads all of the properties represented by the XML document on the
// specified input stream into this properties table.
//
// <p>The XML document must have the following DOCTYPE declaration:
// <pre>
// &lt;!DOCTYPE properties SYSTEM "http://java.sun.com/dtd/properties.dtd"&gt;
// </pre>
// Furthermore, the document must satisfy the properties DTD described
// above.
//
// <p> An implementation is required to read XML documents that use the
// "{@code UTF-8}" or "{@code UTF-16}" encoding. An implementation may
// support additional encodings.
//
// <p>The specified stream is closed after this method returns.
//
// @param in the input stream from which to read the XML document.
// @throws IOException if reading from the specified input stream
//         results in an <tt>IOException</tt>.
// @throws java.io.UnsupportedEncodingException if the document's encoding
//         declaration can be read and it specifies an encoding that is not
//         supported
// @throws InvalidPropertiesFormatException Data on input stream does not
//         constitute a valid XML document with the mandated document type.
// @throws NullPointerException if {@code in} is null.
// @see    #storeToXML(OutputStream, String, String)
// @see    <a href="http://www.w3.org/TR/REC-xml/#charencoding">Character
//         Encoding in Entities</a>
// @since 1.5
func (p *Properties) LoadFromXML(reader io.Reader) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return load(p, reader)
}

// Emits an XML document representing all of the properties contained
// in this table.
//
// <p> An invocation of this method of the form <tt>props.storeToXML(os,
// comment)</tt> behaves in exactly the same way as the invocation
// <tt>props.storeToXML(os, comment, "UTF-8");</tt>.
//
// @param os the output stream on which to emit the XML document.
// @param comment a description of the property list, or {@code null}
//        if no comment is desired.
// @throws IOException if writing to the specified output stream
//         results in an <tt>IOException</tt>.
// @throws NullPointerException if {@code os} is null.
// @throws ClassCastException  if this {@code Properties} object
//         contains any keys or values that are not
//         {@code Strings}.
// @see    #loadFromXML(InputStream)
// @since 1.5
func (p *Properties) StoreToXML(writer io.Writer, comments string) error {
	return p.StoreToXMLByEncoding(writer, comments, "UTF-8")
}

// Emits an XML document representing all of the properties contained
// in this table, using the specified encoding.
//
// <p>The XML document will have the following DOCTYPE declaration:
// <pre>
// &lt;!DOCTYPE properties SYSTEM "http://java.sun.com/dtd/properties.dtd"&gt;
// </pre>
//
// <p>If the specified comment is {@code null} then no comment
// will be stored in the document.
//
// <p> An implementation is required to support writing of XML documents
// that use the "{@code UTF-8}" or "{@code UTF-16}" encoding. An
// implementation may support additional encodings.
//
// <p>The specified stream remains open after this method returns.
//
// @param os        the output stream on which to emit the XML document.
// @param comment   a description of the property list, or {@code null}
//                  if no comment is desired.
// @param  encoding the name of a supported
//                  <a href="../lang/package-summary.html#charenc">
//                  character encoding</a>
//
// @throws IOException if writing to the specified output stream
//         results in an <tt>IOException</tt>.
// @throws java.io.UnsupportedEncodingException if the encoding is not
//         supported by the implementation.
// @throws NullPointerException if {@code os} is {@code null},
//         or if {@code encoding} is {@code null}.
// @throws ClassCastException  if this {@code Properties} object
//         contains any keys or values that are not
//         {@code Strings}.
// @see    #loadFromXML(InputStream)
// @see    <a href="http://www.w3.org/TR/REC-xml/#charencoding">Character
//         Encoding in Entities</a>
// @since 1.5
func (p *Properties) StoreToXMLByEncoding(writer io.Writer, comments, encoding string) error {
	return save(p, writer, comments, encoding)
}

// Searches for the property with the specified key in this property list.
// If the key is not found in this property list, the default property list,
// and its defaults, recursively, are then checked. The method returns
// {@code null} if the property is not found.
//
// @param   key   the property key.
// @return  the value in this property list with the specified key value.
// @see     #setProperty
// @see     #defaults
func (p *Properties) GetProperty(key string) (string, error) {
	var oval = p.Get(key)
	if sval,ok := oval.(string); ok {
		return sval, nil
	}

	if p.defaults != nil {
		return p.defaults.GetProperty(key)
	}

	return "", errors.New("<nil>")
}

// Searches for the property with the specified key in this property list.
// If the key is not found in this property list, the default property list,
// and its defaults, recursively, are then checked. The method returns the
// default value argument if the property is not found.
//
// @param   key            the hashtable key.
// @param   defaultValue   a default value.
//
// @return  the value in this property list with the specified key value.
// @see     #setProperty
// @see     #defaults
func (p *Properties) GetPropertyByDefault(key, defaultValue string) string {
	value,err := p.GetProperty(key)
	if err != nil {
		return defaultValue
	}
	return value
}

// Returns an enumeration of all the keys in this property list,
// including distinct keys in the default property list if a key
// of the same name has not already been found from the main
// properties list.
//
// @return  an enumeration of all the keys in this property list, including
//          the keys in the default property list.
// @throws  ClassCastException if any key in this property list
//          is not a string.
// @see     java.util.Enumeration
// @see     java.util.Properties#defaults
// @see     #stringPropertyNames
func (p *Properties) PropertyNames() []interface{} {
	var h = p.NewHashtable()
	h.Init()
	p.enumerate(h)

	return h.Keys()
}

// Returns a set of keys in this property list where
// the key and its corresponding value are strings,
// including distinct keys in the default property list if a key
// of the same name has not already been found from the main
// properties list.  Properties whose key or value is not
// of type <tt>String</tt> are omitted.
// <p>
// The returned set is not backed by the <tt>Properties</tt> object.
// Changes to this <tt>Properties</tt> are not reflected in the set,
// or vice versa.
//
// @return  a set of keys in this property list where
//          the key and its corresponding value are strings,
//          including the keys in the default property list.
// @see     java.util.Properties#defaults
// @since   1.6
func (p *Properties) StringPropertyNames() []string {
	var h = p.NewHashtable()
	h.Init()
	p.enumerateStringProperties(h)

	var set []string
	for _,key := range h.Keys() {
		set = append(set, key.(string))
	}

	return set
}

// Rather than use an anonymous inner class to share common code, this
// method is duplicated in order to ensure that a non-1.1 compiler can
// compile this file.
func (p *Properties) List(out io.Writer) {
	out.Write(append([]byte("-- listing properties --"), newLine()...))

	var h = p.NewHashtable()
	h.Init()
	p.enumerate(h)
	for _,key := range h.Keys() {
		var val = h.Get(key)
		var skey = key.(string)
		var sval = val.(string)

		if len(sval) > 40 {
			sval = string([]byte(sval)[:37]) + "..."
		}
		out.Write(append([]byte(skey + " = " + sval), newLine()...))
	}
}

// Enumerates all key/value pairs in the specified hashtable.
// @param h the hashtable
// @throws ClassCastException if any of the property keys
//         is not of String type.
func (p *Properties) enumerate(h Hashtable) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.defaults != nil {
		p.defaults.enumerate(h)
	}

	for _,key := range p.Hashtable.Keys() {
		h.Put(key.(string), p.Hashtable.Get(key))
	}
}

// Enumerates all key/value pairs in the specified hashtable
// and omits the property if the key or value is not a string.
// @param h the hashtable
func (p *Properties) enumerateStringProperties(h Hashtable) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.defaults != nil {
		p.defaults.enumerateStringProperties(h)
	}

	for _,key := range p.Hashtable.Keys() {
		var val = p.Hashtable.Get(key)
		if _,ok := key.(string); ok {
			if _,ok = val.(string); ok {
				h.Put(key.(string), val.(string))
			}
		}
	}
}

// Create a hashtable that is the same as itself
func (p *Properties) NewHashtable() Hashtable {
	t := reflect.TypeOf(p.Hashtable).Elem()
	v := reflect.New(t)
	return v.Interface().(Hashtable)
}

func (p *Properties) ToMap() map[interface{}]interface{} {
	var m = make(map[interface{}]interface{})
	keys := p.Keys()
	for _,key := range keys {
		m[key] = p.Get(key)
	}

	return m
}

func (p *Properties) New() *Properties {
	t := reflect.TypeOf(*p)
	v := reflect.New(t)
	var pro = v.Interface().(*Properties)
	pro.Hashtable = p.NewHashtable()
	pro.Init()

	return pro
}


// Convert a nibble to a hex character
// @param   nibble  the nibble to convert.
func toHex(nibble int) byte {
	return hexDigit[nibble & 0xF]
}

// Write a comments.
func writeComments(bw *bufio.Writer, comments string) (err error) {
	if err = bw.WriteByte('#'); err != nil {
		return err
	}
	var length = len(comments)
	var current = 0
	var last = 0
	var uu = make([]byte, 6)
	uu[0] = '\\'
	uu[1] = 'u'
	for current < length {
		var c = comments[current]
		if c > '\u00ff' || c == '\n' || c == '\r' {
			if last != current {
				if _,err = bw.Write([]byte(comments)[last:current]); err != nil {
					return err
				}
			}
			if c > '\u00ff' {
				uu[2] = toHex(int(c >> 12) & 0xf)
				uu[3] = toHex(int(c >>  8) & 0xf)
				uu[4] = toHex(int(c >>  4) & 0xf)
				uu[5] = toHex(int(c)       & 0xf)
				if _,err = bw.Write(uu); err != nil {
					return err
				}
			} else {
				if _,err = bw.Write(newLine()); err != nil {
					return err
				}
				if c == '\r' && current != length - 1 && comments[current + 1] == '\n' {
					current++
				}
				if current == length - 1 || comments[current + 1] != '#' && comments[current + 1] != '!' {
					if err = bw.WriteByte('#'); err != nil {
						return err
					}
				}
			}
			last = current + 1
		}
		current++
	}

	if last != current {
		if _,err = bw.Write([]byte(comments)[last:current]); err != nil {
			return err
		}
	}
	if _,err = bw.Write(newLine()); err != nil {
		return err
	}

	return nil
}

// Writes a line separator.  The line separator string is defined by the
// system property <tt>line.separator</tt>, and is not necessarily a single
// newline ('\n') character.
func newLine() []byte {
	const (
		CR = '\r'
		LF = '\n'
	)
	switch runtime.GOOS {
	case "windows":
		return []byte{CR, LF}
	case "linux":
		fallthrough
	default:
		return []byte{LF}
	}

	return []byte{LF}
}



// Read in a "logical line" from an InputStream/Reader, skip all comment
// and blank lines and filter out those leading whitespace characters
// (\u0020, \u0009 and \u000c) from the beginning of a "natural line".
// Method returns the char length of the "logical line" and stores
// the line in "lineBuf".
type LineReader interface {
	readLine() int
}

type lineReader struct {
	inByteBuf    []byte
	lineBuf      []byte

	inLimit       int
	inOff         int

	reader        io.Reader
}

func NewLineReader(reader io.Reader) LineReader {
	return &lineReader{
		inByteBuf : make([]byte, 8192),
		lineBuf   : make([]byte, 1024),
		reader    : reader,
		inLimit   : 0,
		inOff     : 0,
	}
}

func (l *lineReader) readLine() int {
	var length = 0
	var c byte = 0
	var (
		skipWhiteSpace     = true
		isCommentLine      = false
		isNewLine          = true
		appendedLineBegin  = false
		precedingBackslash = false
		skipLF             = false
	)

	for true {
		if l.inOff >= l.inLimit {
			n,err := l.reader.Read(l.inByteBuf)
			l.inLimit = n
			l.inOff = 0
			if err != nil || l.inLimit <= 0 {
				if length == 0 || isCommentLine {
					return -1
				}
				if precedingBackslash {
					length--
				}
				return length
			}
		}

		c = l.inByteBuf[l.inOff]
		l.inOff++

		if skipLF {
			skipLF = false
			if c == '\n' {
				continue
			}
		}
		if skipWhiteSpace {
			if c == ' ' || c == '\t' || c == '\f' {
				continue
			}
			if !appendedLineBegin && (c == '\r' || c == '\n') {
				continue
			}
			skipWhiteSpace = false
			appendedLineBegin = false
		}
		if isNewLine {
			isNewLine = false
			if c == '#' || c == '!' {
				isCommentLine = true
				continue
			}
		}

		if c != '\n' && c != '\r' {
			l.lineBuf[length] = c
			length++
			if length == len(l.lineBuf) {
				var newLength = length * 2
				if newLength < 0 {
					newLength = int(^uint(0) >> 1)
				}
				var buf = make([]byte, newLength)
				copy(buf, l.lineBuf)
				l.lineBuf = buf
			}
			//flip the preceding backslash flag
			if c == '\\' {
				precedingBackslash = !precedingBackslash
			} else {
				precedingBackslash = false
			}
		} else {
			// reached EOL
			if isCommentLine || length == 0 {
				isCommentLine = false
				isNewLine = true
				skipWhiteSpace = true
				length = 0
				continue
			}
			if l.inOff >= l.inLimit {
				n,err := l.reader.Read(l.inByteBuf)
				l.inLimit = n
				l.inOff = 0
				if err != nil || l.inLimit <= 0 {
					if precedingBackslash {
						length--
					}
					return length
				}
				if precedingBackslash {
					length -= 1
					//skip the leading whitespace characters in following line
					skipWhiteSpace = true
					appendedLineBegin = true
					precedingBackslash = false
					if c == '\r' {
						skipLF = true
					}
				} else {
					return length
				}
			}
			if precedingBackslash {
				length -= 1
				//skip the leading whitespace characters in following line
				skipWhiteSpace = true
				appendedLineBegin = true
				precedingBackslash = false
				if c == '\r' {
					skipLF = true
				}
			} else {
				return length
			}
		}
	}

	return -1
}
