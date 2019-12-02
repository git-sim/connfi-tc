# Messaging Service TC (UNDER CONSTRUCTION)

> This messaging service is composed of two parts (backend & frontend) running in their own containers.

## Quick Start
The environment params (Ports and options) are setup in the .env file. Use docker-compose to start up both containers. 
Download this github repo, cd into the directory containing the toplevel docker-compose.yml and call
``` bash
docker-compose up
```
Then access the frontend client to setup an account and start sending messages. 

# Backend App
> Executes the business rules, persistence, and the REST API for the messaging service. The backend is designed based on "Clean Architecture" principles.  Their are 3 regions separated by boundaries. The domain region contains the business/domain logic. The usecase region defines the types and interactors needed to carry out the usecases of the system.  The IO layer contains the details for talking to the system (Web, Storage)

> The business logic is located in the /domain directory, with the subdirectories:
	>	domain/entity - Contains the business objects that aren't dependent on any components
	>   domain/repo   - Defines interfaces for the repositories providing persistence for the entities 
	>   domain/service - Another layer for dependency inversion so the Entities don't have to know about repo's at all.
> The usecase interactors are in the /usecase directory. The major usecases of system revolve around Messaging, Account management, Profile management.
> The IO details are confined to the /IO directory many of implementations of the interfaces defined in domain/repo are found here.
    >   IO/storage/ram    - implementation for the {Domain | Storage} and {Usecase | Storage} boundaries
        >   IO/storage/ram    - ram based inmemory implementation of the domain repos
	    >   IO/storage/mdb    - mongodb implemenations of the domain repos. Not yet implemented
	>   IO/restapi/v1.0/  - the restapi implementation for the {HTTP | Usecase} boundary.  Not yet implemented
	
# Frontend Client Single Page Application
> Runs the GUI elements, interacts with the user, and communicates with the backend App's container over a REST ifc.
> Not yet implemented


