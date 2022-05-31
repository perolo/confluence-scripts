# confluence-scripts
Handy utilities for administration and maintenance of Confluence
The use of properties files is to enable automation with CI tool without "logging" passwords.

## Clone Group Permissions
Confluence does not support renaming of groups.
This script clones the permissions of a group in all Spaces to a new group.

Does not clone the group members!

##### Requires plugin: REST Extender for Confluence   https://it-lab-site.atlassian.net/wiki/spaces/RAEC/overview 
Version 2.4.x 

#### How to use
* Create the group(copygroupname) in Confluence
* Modify the properties file 
* Run the script
        ```
        go run clonegrouppermissions.go -prop clonegrouppermissions.properties
        ```
* (Optional) Remove the old group 

The parameters needed are defined in the properties file:  
```
confhost=http://confluence.com:8080
user=user
password=password
source=origingroupname
destination=copygroupname
```
## Space Permissions Report
Creates a Excel Sheet with the users and groups and the permissions in a SpaceCategory
Requires a user with admin permissions in Confluence 

* Modify the properties file
* Run the script        ```
                        go run .\main.go -prop .\spacepermissionsreport.properties
                        ```
 
The parameters needed are defined in the properties file:  
```
confhost=http://confluence.com:8080
user=user
password=password
spacecategory=tools
groups=true
users=true
file=C:/Tmp/Report.xlsx
```
##### Requires plugin: REST Extender for Confluence   https://it-lab-site.atlassian.net/wiki/spaces/RAEC/overview 


## Space Permissions Modifier
Inteded as an interactive script that guides you to modify the permissions for users and groups in a SpaceCategory
Requires a user with admin permissions in Confluence 

>:warning: **Works Ok in Linux, but fails in WIndows : https://github.com/manifoldco/promptui/issues/33**

>>[!WARNING]
>Works Ok in Linux, but fails in WIndows : https://github.com/manifoldco/promptui/issues/33

* Modify the properties file:
* Run the script        ```
                        go run spacepermissionsmodifier.go -prop spacepermissionsmodifier.properties
                        ```

The parameters needed are defined in the properties file:  
```
confhost=http://confluence.com:8080
user=user
password=password
spacecategory=tools
groups=true
users=true
```
##### Requires plugin: REST Extender for Confluence   https://it-lab-site.atlassian.net/wiki/spaces/RAEC/overview 

## AddGroupPage
Creates a Confluence page with the confluence User List for all groups in Confluence

* Modify the properties file
* Run the script        ```
                        go run .\addgrouppage\main\main.go -prop .\addgrouppage\main\addgrouppage.properties
                        ```
 
The parameters needed are defined in the properties file:  
```
confhost=http://192.168.50.40:8090
confuser=admin
confpass=admin
conftoken=false
usetoken=false
space=ds
ancestortitle=Welcome to Confluence

```

## Personal Spaces Report
Creates a Excel Sheet with the users and groups and the permissions in all personal spaces.
Could be used to review/archive personal spaces of deactivated users.
Requires a user with admin permissions in Confluence 

* Modify the properties file
* Run the script        ```
                        go run .\main.go -prop .\spacepermissionsreport.properties
                        ```
 
The parameters needed are defined in the properties file:  
```
confhost=http://192.168.50.40:8090
confuser=admin
confpass=admin
conftoken=false
usetoken=false
space=ds
ancestortitle=Welcome to Confluence
report=true
file=C:/temp/Report%s.xlsx
simple=false
```
##### Requires plugin: REST Extender for Confluence   https://it-lab-site.atlassian.net/wiki/spaces/RAEC/overview 



## Sync AD Group
(under development)
Intended as a user/permission synchronization script with a structure defined in the AD.
For tools not directly integrated/connected to AD
## Get GitLab User report
(under development)
Intended as a user/permission review report comparing with a structure defined in the AD
## Get JIRA User report
(under development)
Intended as a user/permission review report comparing with a structure defined in the AD
## Get SonarQube User report
(under development)
Intended as a user/permission review report comparing with a structure defined in the AD
