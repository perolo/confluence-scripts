package htmlutils

import "os"

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func IsWritable(name string) (isWritable bool) {
	isWritable = false

	_, err := os.OpenFile(name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return isWritable
	}
	isWritable = true
	return isWritable
}

func WriteHeader1(w *os.File, s string) {
	_, err := w.WriteString("<h1>" + s + "</h1>\n")
	Check(err)
}

func WriteHeader2(w *os.File, s string) {
	_, err := w.WriteString("<h2>" + s + "</h2>\n")
	Check(err)
}

func WriteParagraf(w *os.File, s string) {
	_, err := w.WriteString("<p>" + s + "</p>\n")
	Check(err)
}
func WriteBullitStart(w *os.File) {
	_, err := w.WriteString("<ul>\n")
	Check(err)
}
func WriteBullit(w *os.File, s string) {
	_, err := w.WriteString("   <li>" + s + "</li>\n")
	Check(err)
}
func WriteBulliEnd(w *os.File) {
	_, err := w.WriteString("</ul>\n")
	Check(err)
}

func WriteTableHeader1(w *os.File, s []string) {

	_, err := w.WriteString("<table><tr>")
	Check(err)
	for _, header := range s {
		_, err2 := w.WriteString("<th>" + header + "</th>")
		Check(err2)
	}
	_, err = w.WriteString("</tr>\n")
	Check(err)
}

func WriteTableLine(w *os.File, s []string) {

	_, err := w.WriteString("<tr>")
	Check(err)
	for _, header := range s {
		_, err2 := w.WriteString("<td>" + header + "</td>")
		Check(err2)
	}
	_, err = w.WriteString("</tr>\n")
	Check(err)
}

func WriteTableEnd(w *os.File) {
	_, err := w.WriteString("</table>\n")
	Check(err)
}
func WriteWrapLink(w *os.File, link string, name string) {
	_, err := w.WriteString(WrapLink(link, name))
	Check(err)
}

func WrapLink(link string, name string) string {
	if name == "" {
		return ""
	}
	return "<a href=" + link + ">" + name + "</a>"
}
func WrapJIRA(key string, hideOptional ...bool) string {
	if key == "" {
		return ""
	}
	hide := false
	if len(hideOptional) > 0 {
		hide = hideOptional[0]
	}
	hideparam := ""
	if hide {
		hideparam = "<ac:parameter ac:name=\"showSummary\">false</ac:parameter>"
	}
	return "<p>	<ac:structured-macro ac:macro-id=\"ef027b1e-9afc-4bf1-bc01-3e464985fb8d\" ac:name=\"jira\" ac:schema-version=\"1\">" +
		hideparam +
		"<ac:parameter ac:name=\"server\">Shared Technologies - JIRA</ac:parameter>" +
		"<ac:parameter ac:name=\"columns\">key,summary,type,created,updated,due,assignee,reporter,priority,status,resolution</ac:parameter>" +
		"<ac:parameter ac:name=\"serverId\">936bba59-626d-360c-aecd-b1292bf65b83</ac:parameter>" +
		"<ac:parameter ac:name=>" + key + "</ac:parameter>" +
		"</ac:structured-macro>	</p>"
}

/*
<p>
<ac:structured-macro ac:macro-id="6f4c01bc-8463-43fc-9434-102afcb931b1" ac:name="contentbylabel" ac:schema-version="3">
<ac:parameter ac:name="cql">label = "comp_c"</ac:parameter>
</ac:structured-macro>
</p>
*/

func WrapLabel(key string) string {
	if key == "" {
		return ""
	}
	return "<p>	<ac:structured-macro ac:macro-id=\"6f4c01bc-8463-43fc-9434-102afcb931b1\" ac:name=\"contentbylabel\" ac:schema-version=\"3\">" +
		"<ac:parameter ac:name=\"cql\">label = \"" + key + "\"</ac:parameter>" +
		"</ac:structured-macro></p>"
}

func WrapBold(text string) string {
	return "<b>" + text + "</b>"
}
