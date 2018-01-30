# Rule specification language
The rules are described via the custom DSL, which syntax is explained in this document alongside with couple [examples](#examples) of the simple rules.

## Rule
The basic syntax of single rule matches the following format:

<pre>
<b>RULE</b> name<b>:</b>
&nbsp;&nbsp;&nbsp;&nbsp;<em><a href=#condition>conditions</a></em>     
<b>TRIGGERS</b>   
&nbsp;&nbsp;&nbsp;&nbsp;<em><a href=#action>actions</a></em>
</pre>

with the following semantics:

|name|description|format|
|:--:|:----------|:----:|
**RULE**|Required keyword for the start of the rule.|
*name*|Arbitrary name of the rule.|ID
**:**|Required syntax sugar.|
*conditions*|List of conditions.|[condition](#condition)|
**TRIGGERS**|Required keyword for start of the actions section.|
*actions*|List of actions to perform after positive evaluation of the conditions.|[action](#action)

Aforementioned sections of the rule must come in that exact order. 

## Condition
Single condition consists of the following sections:

<pre>
<a href="#parameter">Parameter</a>&nbsp;&nbsp;&nbsp;<a href="#operator">Operator</a>&nbsp;&nbsp;&nbsp; <a href="#value">Value</a>
</pre>

For example:
```
8837ffdf-2bec-42f7-9c2d-b8cfa67661a9["temperature"] >= 30
```
To specify multiple conditions for the same rule, just add them one after another. If there are more than one condition specified, logical *AND* will be considered as their binding.

### Parameter
Represents the source of the condition. It has the following syntax:
<pre>
<em>deviceId</em><b>[</b><em>deviceProperty</em><b>]</b>
</pre>

with the following semantics:

|name|description|format|
|:--:|:----------|:----:|
|*deviceId*|Identifier of the device.| UUID
|**[**|Required syntax sugar.|
|*deviceProperty*|Property of the device.|String
|**]**|Required syntax sugar.|

### Operator
Used for comparison of the parameters and value. Supported operators are:

|keyword|meaning|
|:-----:|:------|
|**=**|Equals|
|**!=**|Not equals|
|**<**|Less than|
|**<=**|Less than or equals|
|**>**|Greater than|
|**=>**|Greater than or equals|
|**BETWEEN**|Greater than *and* less than|

### Value
Represents value which is compared to specified device's property.

Supported values:
 - String
 - Integer
 - Float
 - [Range](#range)

#### Range
Range is custom value with the following syntax:

<pre>
<b>[</b><em>lowerBound</em> <b>,</b> <em>upperBound</em><b>]</b>
</pre>

|name|description|format|
|:--:|:----------|:----:|
|**[**|Required syntax sugar.|
|*lowerBound*|Lower bound of the range.|Float
|**,**|Required syntax sugar.|
|*upperBound*|Upper bound of the range.|Float
|**]**|Required syntax sugar.|

## Action
Currently, there are several different action supported
- [Send Email](#send-email-action)
- [Turn Off](#turn-off-action)

To set multiple actions to trigger after positive evaluation of the rule, simply specify them one after another.

### Send Email Action
This action executes email sending to the specified recipient.

<pre>
<b>SEND EMAIL</b> <em>content</em> <b>TO</b> <em>recipient</em>
</pre>

|name|description|format|
|:--:|:----------|:----:|
|**SEND EMAIL**|Keyword that specifies the action name.|
|*content*|The content of the email.|String
|**TO**|Required syntax sugar|
|*recipient*|Email address of the recipient|String

### Turn Off Action
This action executes shutdown signal on the device with the specified identifier.

<pre>
<b>TURN OFF</b> <em>deviceId</em>
</pre>

|name|meaning|format|
|:--:|:------|:----:|
|**TURN OFF**|Keyword that specifies the action name.|
|*deviceId*|Identifier of the device to execute action on.|UUID

## Examples
<pre>
<b>RULE</b> rule01<b>:</b>
    8837ffdf-2bec-42f7-9c2d-b8cfa67661a9<b>[</b>"temperature"<b>]</b> <b>>=</b> 30 
<b>TRIGGERS</b>
    <b>TURN OFF</b> 45ab8a94-e0f1-4ea4-b99c-c51e16566f83
    <b>SEND EMAIL</b> "Alert! High temperature in the living room!" <b>TO</b> "person01@home.com"

<b>RULE</b> rule02<b>:</b>
    5b49e289-9be5-4622-9550-d52bbcbea37a<b>[</b>"temperature"<b>]</b> <b>>=</b> 30
    a772aba0-0065-4af6-bd47-3e956d569f5a<b>[</b>"activeHeatersNumber"<b>]</b> <b>BETWEEN</b> <b>[</b>1, 4<b>]</b>
<b>TRIGGERS</b>
    <b>TURN OFF</b> 23cdb6f4-f785-4c7a-89c7-868ca7100dc5
    <b>TURN OFF</b> 9c4feb41-a705-4b1e-a6fc-7d5e7f1f964e
</pre>

