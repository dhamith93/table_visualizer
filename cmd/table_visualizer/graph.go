package main

func generateGraphForAll(tables []Table) string {
	out := graphStart()
	for _, t := range tables {
		out += makeGraph(&t, true)
	}
	out += "\n}"
	return out
}

func generateGraphForOne(t *Table) string {
	out := graphStart()
	out += makeGraph(t, false)
	out += "\n}"
	return out
}

func makeGraph(t *Table, isAll bool) string {
	out := createTableNode(t.Name)
	for _, c := range t.Columns {
		out += "\n" + t.Name + " -> \"" + t.Name + "." + c.Name + "\""
		for _, fk := range c.Fks {
			out += createTableNode(fk.Table)
			out += "\n\"" + t.Name + "." + c.Name + "\" -> " + "{ \"" + fk.Table + "." + fk.RefCol + "\" } [ style=dashed ]; \"" + fk.Table + "." + fk.RefCol + "\" -> " + fk.Table
			if isAll {
				out += "[style=invis];\n"
			}
		}
	}
	return out + "\n"
}

func graphStart() string {
	return `digraph G {
		graph [
			rankdir=LR
		];
		
		graph [
			splines=polyline
			rankdir=LR
		];
                
		edge [
			arrowhead=normal,
			weight=1
		];`
}

func createTableNode(name string) string {
	return "\n" + name + `[
		height=0.84444,
		margin=0.3,
		shape=box];`
}
