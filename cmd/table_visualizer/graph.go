package main

func generateGraph(tables []Table) string {
	out := graphStart()
	for _, t := range tables {
		out += createGraphTable(&t)
	}
	for _, t := range tables {
		out += createFkRelationships(&t)
	}
	out += "\n}"
	return out
}

func graphStart() string {
	return `digraph G {		
		node [shape=plaintext]
		`
}

func createGraphTable(t *Table) string {
	out := "\n" + t.Name + `[label=<
	<TABLE BORDER="0" CELLBORDER="1" CELLSPACING="0">
	 <TR>
	  <TD BGCOLOR="gray" ALIGN="LEFT">` + t.Name + `     </TD>
	 </TR>`
	for _, c := range t.Columns {
		out += `<TR>
		<TD PORT="` + t.Name + "." + c.Name + `" ALIGN="LEFT">` + c.Name + ` ` + c.Type + `  </TD>
	   </TR>`
	}
	return out + "\n</TABLE>>];\n"
}

func createFkRelationships(t *Table) string {
	out := "\n"
	for _, c := range t.Columns {
		for _, fk := range c.Fks {
			out += fk.Table + ":" + "\"" + fk.Table + "." + fk.RefCol + "\" -> " + t.Name + ":\"" + t.Name + "." + c.Name + "\" [ style=dashed ];\n"
		}
	}
	return out
}
