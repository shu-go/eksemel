eksemel add --xpath \"//command[@name='add']//optiona[@name='name']\" --name #comment --value "ADD" help_wip.xml | ^
eksemel add --xpath \"//command[@name='add']/options\" --ennet "\"option[name=ennet{emmet-like abbreviation}\"" | ^
eksemel delete --xpath \"//command[@name='replace']//optionb[@name='dummy']\" | ^
eksemel replace --xpath \"/zml\" --value eksemel > help.xml
