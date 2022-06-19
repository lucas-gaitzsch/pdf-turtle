package templateengines

import "encoding/json"

const resultHtml = `
<html>
<body>
	<h1>Profile of Bruno</h1>
	<p>Working at Testcompany</p>
	<p>Locations:</p>
	<ul>
		<li>Chemnitz</li><li>Berlin</li><li>Amsterdam</li>
	</ul>
</body>
</html>
`
const jsonModel = `
{
	"name": "Bruno",
	"company": {
		"name": "Testcompany",
		"locations": ["Chemnitz", "Berlin", "Amsterdam"]
	}
}
`

func getModel() any {
	var model any
	json.Unmarshal([]byte(jsonModel), &model)
	return model
}
