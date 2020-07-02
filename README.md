# confluence-scripts
Handy utilities for administration and maintennance of Confluence

## Clone Group Permissions
Confluence does not support renaming of groups.
This script clones the permissions of a group in all Spaces to a new group.

Does not clone the group members!

##### Requires plugin: REST Extender for Confluence   https://it-lab-site.atlassian.net/wiki/spaces/RAEC/overview 
Version 2.4.x - (not yet released?!?!)

#### How to use
* Create the group in Confluence
* Modify the properties file with:
    * Confluence server name
    * User
    * Password
    * Source group in Confluence
    * Destination group in Confluence  
* Run the script
        ```go
        go run clonegrouppermissions.go -prop clonegrouppermissions.properties
        ```
* (Optional) Remove the old group 

The parameters needed are defined in the properties file:  
```go
confhost=http://confluence.com:8080
user=user
password=password
source=origin
destination=copy
```
