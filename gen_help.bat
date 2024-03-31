eksemel add --xpath \"//command[@name='add']//option[@name='name']\" --name #comment --value "ADD" help_wip.xml | ^
eksemel add --xpath \"//command[@name='add']//option[@name='name']\" --name desc --value "element, @attribute, '#text', '#cdata-section', '#comment'" | ^
eksemel add --xpath \"//command[@name='add']//option[@name='name']\" --name required/@required --value true | ^
eksemel add --xpath \"//command[@name='add']//option[@name='value']\" --name #comment --value "ADD" | ^
eksemel add --xpath \"//command[@name='add']//option[@name='value']\" --name desc --value "element also can have --value, to add <name>value</name>" | ^
eksemel add --xpath \"//command[@name='replace']\" --name #comment --value "DEL dummy" | ^
eksemel add --xpath \"//command[@name='replace']\" --name #comment --value "ADD delete" | ^
eksemel delete --xpath \"//command[@name='replace']//option[@name='dummy']\" | ^
eksemel add --xpath \"//command[@name='replace']\" --sibling --name "command/@name" --value "delete" | ^
eksemel add --xpath \"//command[@name='delete']\" --name "options/option/@name" --value "xpath" | ^
eksemel add --xpath \"//command[@name='delete']//option[@name='xpath']\" --name "desc" --value "An XPath to the target" | ^
eksemel replace --xpath \"/xml\" --value eksemel > help.xml
