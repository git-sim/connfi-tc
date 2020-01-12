# Messaging Service TC (UNDER CONSTRUCTION)

> This messaging service is composed of two parts (backend & frontend) running in their own containers.

## Quick Start
The environment params (Ports and options) are setup in the .env file. Use docker-compose to start up both containers. 
Download this github repo, cd into the directory containing the toplevel docker-compose.yml and call
``` bash
sudo docker-compose up -d
```
Then access the frontend client to setup an account and start sending messages. 

## Backend App
Executes the business rules, persistence, and the REST API for the messaging service. The backend is designed based on "Clean Architecture" principles.  
There are 3 regions separated by boundaries. 
* The domain region contains the business/domain logic. 
* The usecase region defines the types and interactors needed to carry out the usecases of the system.  
* The IO region contains the details for talking to the system (Web, Storage)

### Directory structure
* [/domain]()     contains the business logic is located in the /domain directory, with the subdirectories:
  * [./entity]()   Contains the business objects that aren't dependent on any components
  * [./repo]()     Defines interfaces for the repositories providing persistence for the entities 
  * [./service]()   A layer for dependency inversion for the usecases so the Entities don't have to know about usecase logic.

* [/usecase]() contains usecase interactors 
    * The major usecases of system revolve around Messaging, Account management, Profile management.
    * Add more details about the usecases...could be its own section.

* [/IO]() contains the IO details of the system. Implementations of the interfaces defined in domain/repo are found here.
  * [./storage]()  implementation for the {Domain | Storage} and {Usecase | Storage} boundaries
    * [./ram]()    ram based inmemory implementation of the domain repos for testing and demos
    * [./mdb]()    mongodb implementations of the domain repos. Not yet implemented
  * [./rest]()  the restapi implementation for the {HTTP | Usecase} boundary.
    * The endpoints are 

      * [/login?email=val]() 
        Method * - Create a new login session
      * [/logout/{accountID}]
        Method * - Logout a session

### REST API

#### accounts: Information about user accounts. Contains id, email, firstname, lastname  
* POST /accounts
  * Registers a New Account in the system   
  * Input: Account as body param   
    * {email,firstname,lastname}   
  * Output: Returns id of the newly created account if email was unique   

* GET /accounts
  * Get the list of accounts.  
  * Input:  limit and offset are optional. If not specified means all accounts.   
    * ?[limit=n]  
    * &[offset=n]  
  * Output: {TotalNumOfAccounts, Accounts[]}  

* GET /accounts/{accountID}
  * Returns specific account info.  
  * Input: none  
  * Output: {ID, email, firstname, lastname}   

* PUT /accounts/{accountID}
  * Updates the account info   
  * Input: AccountInfo {email,firstname,lastname}   
  * Output: none

* DELETE /accounts/{accountID}
  * Delete account 
  * Input: none
  * Output: none

#### folders: Contains summary info about user folders (inbox, archive, etc). A particular folder can be queried for all the messages in the folder, with sorting. 
* GET /accounts/{accountID}/folders
  * Returns list of folders  
  * Input: none  
  * Output: {TotalNumberOfFolders, FolderInfo[]}  
    * Where FolderInfo:= {FolderName, Idx, NumTotal, NumUnviewed} 

* GET /accounts/{accountID}/folders/{folderID}  
  * Returns the messages in a folder sorted/limited/paged for the frontend. Page size is specified by limit. So {Limit:10,Page:0} gives the first 10 messages.  {Limit:10,Page:1} gives the next 10. 
  * Input:  
    * ?[limit=n]  
    * &[page=n]  
    * &[sortorder=**-1**\|1]  
    * &[sort=**time**\|sender\|subject]  
  * Output: {HeaderInfo, Messages[]}  
    * Where HeaderInfo is {{InputQueryParams}, FolderInfo} 

#### messages: Access to messages for a particular user regardless of what folder they are in. Mainly used for creating, deleting, and marking messages as read. Retrieval/Display of messages is best done via the 'account/{accountID}/folders/{folderID}' endpoint. 

* POST /accounts/{accountID}/messages  
  * Creates a new message  
  * Input:  message as body param  
  * Output: Returns message id of created message  

* GET /accounts/{accountID}/messages  
  * List of messages limit and offset are optional. If not specified means all messages.    
  * Input:  
    * [?limit=n&offset=n]  
  * Output: Returns the total number of messages, and a list of messages(limit,offset)  
            {TotalNumberOfMessages,  Messages[]}  

* GET /accounts/{accountID}/messages/{messageID}  
  * Returns a specific message  
  * Input:  none  
  * Output: {Message}  

* PUT /accounts/{accountID}/messages/{messageID}  
  * Modify a message in a folder to mark it as viewed.  
  * Input:  viewed=0\|1  
  * Output: none  

* DELETE /accounts/{accountID}/messages/{messageID}  
  * Deletes a message  
  * Input:  none  
  * Output: none  

#### threads: Used for displaying messages as threads. Possibly move to under folders.
* GET /accounts/{accountID}/threads
  * Returns list of threadInfos for the account    
  * Input: limit and offset are optional, means all if omitted  
    * [?limit=n&offset=n]  
  * Output: {TotalNumberOfThreads, numRet, ThreadInfo[]}  
    * where ThreadInfo := {IsMuted, NumTotalMsgsInThread, NumUnvieiwedMsgsInThread, MessageID[]}   
* GET /accounts/{accountID}/threads/{threadID}
  * Returns info for a particular thread  
  * Input: none   
  * Output: {ThreadInfo}  (see above)  
* PUT /accounts/{accountID}/threads/{threadID} 
  * Sets the mute param for the thread   
  * Input: mute=0|1    
  * Output: none   
* DELETE /accounts/{accountID}/threads/{threadID}
  * Deletes all the messages in the thread   
 
#### profile: Get/Set the profile for a user.  Broken out by fields to encourage use of the individual requrests. Since the profile may get data/added removed keeps the API the the least coupled if clients ask for the profile fields by name instead of the whole profile. 
* GET /account/{accountID}/profile  DEPRECATED  use the individual fields in the url  
  * Returns the entire profile  
  * Input: none  
  * Output: {Profile}  
    * where Profile:= {ID, name, bio, avatarImg, bgImg }  
* GET /account/{accountID}/profile/name  
  * Returns name string  
  * Input: none  
  * Output: {name}  
* GET /account/{accountID}/profile/bio  
  * Returns bio string  
  * Input: none  
  * Output: {bio}  
* GET /account/{accountID}/profile/avatar  
  * Returns avatar image as a .png  
  * Input: none  
  * Output: {avatar}  
* GET /account/{accountID}/profile/bg  
  * Returns background image as a .png  
  * Input: none  
  * Output: {bg}  
PUT's should be symmetric to the gets above 
* PUT /account/{accountID}/profile  DEPRECATED  use the individual fields in the url  
  * Sets the entire profile    
  * Input: {Profile}  
    * where Profile:= {ID, name, bio, avatarImg, bgImg }  
  * Output:  none  
* PUT /account/{accountID}/profile/name   
  * Sets name string  
  * Input: {name}  
  * Output: none   
* PUT /account/{accountID}/profile/bio  
  * Sets bio string  
  * Input: {bio}  
  * Output: none  
* PUT /account/{accountID}/profile/avatar  
  * Sets avatar image as a .png  
  * Input: none  
  * Output: {avatar}  
* PUT /account/{accountID}/profile/bg  
  * Sets background image as a .png  
  * Input: {bg}    
  * Output: none  

  ---- Original API here for comparison Remove when the refactored api is live ----

    * The endpoints are 
      * [localhost:8080/login?email=val]()  
        * Logs in / Registers a new user
      * [localhost:8080/logout?accid=val]() 
        * Logs out the session
      * [localhost:8080/account?email=val]() 
        * CRU No delete (hotel california). 
        * Most operations require
      * [localhost:8080/accountList]() 
        * Returns the directory info, for autocomplete on FE
        * {Email:"val", ID:"val 64bit hexstring", FirstName:"name", LastName:"name"}
      * [localhost:8080/profile?accid=<val>]() 
        * Not implemented there is a basic CRUD functionality for 
        * Name, Bio, Avatar Image, Background Image. 
        * Plumbed through but not tested at all.
      * [localhost:8080/folder?accid=<val>&msgid=<val>]() 
        * This is used to retrieve a sorted set of messages from a folder (ie inbox)        
        * Optional params: 
      * [localhost:8080/message?accid=<val>&msgid=<val>]()
        * A POST enters a new message into the system for delivery (including scheduled messages).  
        * If a recipient email isn't registered the message is queued up in a pending repo
        * Whenever a CreateUserEvent fires a Listener reads the pending queue gathers any messages for the new user. 

	
## Frontend Client Single Page Application 
Runs the GUI elements, interacts with the user, and communicates with the backend App's container over a REST ifc.
The major components are:
  * Top - The top level component that inherits from React.Component
  * Folders - Select which folder to display messages from. Retrieves displays counters. 
  * Messages - Retrieves and Displays Paginated list of messages from the currently selected folder
  * ComposeForm - Form for composing a new message
  * MessageView - Displays the full body of the actively selected message
  * CreateMessage - Used for generating a Reply/ReplyAll (needs to be combined with ComposeForm)

![Front End Screenshot](https://github.com/git-sim/tc/blob/master/fe_screenshot.PNG)



