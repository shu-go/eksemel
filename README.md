An XML file manipulator

[![Go Report Card](https://goreportcard.com/badge/github.com/shu-go/eksemel)](https://goreportcard.com/report/github.com/shu-go/eksemel)
![MIT License](https://img.shields.io/badge/License-MIT-blue)

# eksemel

Add, update and delete from an XML file by XPath.

## An example (Windows batchfile)

```xml
<?xml version="1.0" standalone="yes"?>
<xml>
    <command name="add">
        <options>
            <option name="xpath">
                <required required="true" />
                <desc>An XPath to the parent</desc>
            </option>
            <option name="name">
                <desc>The name of added node</desc>
            </option>
            <option name="value">
                <desc>@attrname, '#cdata-section', '#text' and '#comment' can have --value</desc>
            </option>
        </options>
    </command>
    <command name="replace">
        <options>
            <option name="xpath">
                <required required="true" />
                <desc>An XPath to the target</desc>
            </option>
            <option name="value">
                <desc>new value</desc>
            </option>
            <option name="dummy">
                <desc>this is a false</desc>
            </option>
            <synopsis>eksemel replace --xpath "//a/b/c/text()" --value "new text"</synopsis>
        </options>
    </command>
</xml>
```

```bat
eksemel add --xpath \"//command[@name='add']//option[@name='name']\" --name #comment --value "ADD" help_wip.xml | ^
eksemel add --xpath \"//command[@name='add']/options\" --ennet "\"option[name=ennet]{emmet-like abbreviation}\"" | ^
eksemel delete --xpath \"//command[@name='replace']//option[@name='dummy']\" | ^
eksemel replace --xpath \"/xml\" --value eksemel > help.xml

eksemel get --xpath \"//desc/text()\" --multiple help.xml
```

# Install

## GitHub Releases

https://github.com/shu-go/eksemel/releases

## Go install

```sh
go install github.com/shu-go/eksemel@latest
```
