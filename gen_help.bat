eksemel add --xpath \"//command[@name='add']//option[@name='name']\" --name #comment --value "ADD" help_wip.xml | ^
eksemel add --xpath \"//command[@name='add']//option[@name='name']\" --name desc --value "element, @attribute, '#text', '#cdata-section', '#comment'" | ^
eksemel add --xpath \"//command[@name='add']//option[@name='value']\" --name #comment --value "ADD" | ^
eksemel add --xpath \"//command[@name='add']//option[@name='value']\" --name desc --value "element also can have --value, to add <name>value</name>" | ^
eksemel add --xpath \"//command[@name='add']/options\" --ennet "\"option[name=ennet]^>desc{emmet-like abbreviation}\"" | ^
eksemel add --xpath \"//command[@name='replace']\" --name #comment --value "DEL dummy" | ^
eksemel add --xpath \"//command[@name='replace']\" --name #comment --value "ADD delete" | ^
eksemel delete --xpath \"//command[@name='replace']//option[@name='dummy']\" | ^
eksemel replace --xpath \"//command[@name='replace']//option[@name='replaced']\" --ennet "\"option[name=ennet]^>desc{emmet-like abbreviation}\"" | ^
eksemel add --xpath \"//command[@name='replace']\" --sibling --ennet \"command[name=delete]\" | ^
eksemel add --xpath \"//command[@name='delete']\" --ennet \"options>option[name=xpath]\"  | ^
eksemel add --xpath \"//command[@name='delete']//option[@name='xpath']\" --name "desc" --value "An XPath to the target" | ^
eksemel add --xpath \"//command[@name='replace']\" --sibling --name #comment --value "vvv by ennet vvv" | ^
eksemel add --xpath \"//command[@name='delete']\" --sibling --ennet "\"command[name=delete]^>options^>option[name=xpath]^>desc{An XPath to the target}\"" | ^
eksemel replace --xpath \"/xml\" --value eksemel > help.xml
