package handlers

import (
    "fmt"
    "net/http"
    "github.com/git-sim/tc/app/usecase"
    "strconv"
)

// The pattern in this code is String based profile fields then Image based profile fields (others if they come alone)
// todo split the two in to separate files?
// Much of the code is very similar, haven't found a way in go to share the code between string and image fields (generics?).
// This is ok since images and strings don't meet the Liskov Substitutability
const (
    EnumFirstNameUsecase = iota
    EnumLastNameUsecase
    EnumBioUsecase
    EnumNumStringUsecases
)

const (
    EnumAvatarImageUsecase = iota
    EnumBgImageUsecase
    EnumNumImageUsecases
)

type ProfileUsecases struct {
    StrUsecases [EnumNumStringUsecases]usecase.ProfileStringUsecase
    ImageUsecases [EnumNumImageUsecases]usecase.ProfileImageUsecase
}

var StringFields = [...]string {
    "firstname",
    "lastname",
    "bio",
}

var ImageFields = [...]string {
    "avatarimage",
    "bgimage",
}

// Helper function to parse the fields in an http request and call the usecase function
func parseAndSetStringUsecase(id64 uint64, fields [EnumNumStringUsecases]string,u *ProfileUsecases, w http.ResponseWriter, r *http.Request) (numParsed int, numErrors int) {
    numParsed = 0
    numErrors = 0
    for i,k := range fields {
        if val, found := r.Form[k]; found {
            err := u.StrUsecases[i].Set(id64,val[0])
            numParsed++
            if err != nil {
                http.Error(w, "err in handleProfile string field",http.StatusBadRequest)
                numErrors++
                // report all errors
            }
        }
    }
    return numParsed, numErrors
}

func parseAndGetStringUsecase(id64 uint64, fields [EnumNumStringUsecases]string,u *ProfileUsecases, w http.ResponseWriter, r *http.Request) (numParsed int, numErrors int) {
    numParsed = 0
    numErrors = 0
    for i,k := range fields {
        if _, found :=  r.Form[k]; found {
            val, err := u.StrUsecases[i].Get(id64)
            numParsed++
            if err != nil {
                http.Error(w, "err in get handleProfile string field",http.StatusBadRequest)
                numErrors++
                // report all errors
            }
            _,_ = fmt.Fprintf(w, "%s\n", val)  //todo json instead
        }
    }
    return numParsed, numErrors
}

//func parseAndSetImageUsecase(id64 uint64, fields [EnumNumImageUsecases]string,u *ProfileUsecases, w http.ResponseWriter, r *http.Request) (numParsed int, numErrors int) {
//    numParsed = 0
//    numErrors = 0
//    for i,k := range fields {
//        val := r.URL.Query().Get(k)
//        if val != nil {
//            numParsed++
//            err := u.ImageUsecases[i].Set(id64,val)
//            if err != nil {
//                http.Error(w, "err in handleProfile image field",http.StatusBadRequest)
//                numErrors++
//                // report all errors
//            }
//        }
//   }
//    return numParsed, numErrors
//}

func HandleProfile(accu usecase.AccountUsecase, u *ProfileUsecases) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        r.ParseForm()
        email := r.FormValue("email")
        if email == "" {
            http.Error(w, "missing email", http.StatusBadRequest)
            return
        }
        account, err := accu.GetAccount(email)
        if err != nil {
            http.Error(w, "email not found", http.StatusBadRequest)
            return
        }
        id64, err := strconv.ParseUint(account.ID,16,64)
        if err != nil {
            http.Error(w, "id lookup failed HandlerProfile", http.StatusInternalServerError)
            return
        }

        switch r.Method {
        case http.MethodPost:
            numStrsParsed,   numStrsErrors   := parseAndSetStringUsecase(id64,StringFields,u,w,r)
            //          numImagesParsed, numImagesErrors := parseAndSetImageUsecase(id64,ImageFields,u,w,r)
            if (numStrsParsed - numStrsErrors > 0) /*|| (numImagesParsed - numImagesErrors > 0)*/  {
                w.WriteHeader(http.StatusCreated)
            }

        case http.MethodDelete:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            // todo allow delete?

        case http.MethodGet:
            numStrsParsed,   numStrsErrors   := parseAndGetStringUsecase(id64,StringFields,u,w,r)
            if(numStrsParsed > 0 && numStrsErrors == 0) {
                w.WriteHeader(http.StatusOK)
            }

        case http.MethodPut:
            numStrsParsed,   numStrsErrors   := parseAndSetStringUsecase(id64,StringFields,u,w,r)
            numImagesParsed := 0 //todo when images are working
            numImagesErrors := 0 //todo images
            //numImagesParsed, numImagesErrors := parseAndSetImageUsecase(id64,ImageFields,u,w,r) //todo images
            if (numStrsErrors + numImagesErrors == 0) && (numStrsParsed + numImagesParsed > 0)  {
                // we parsed at least 1 field, and there were not errors
                w.WriteHeader(http.StatusOK)
            }

        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
        
    })
}

func HandleProfileList(accu usecase.AccountUsecase, u *ProfileUsecases) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch r.Method { 
        case http.MethodGet:
            accs, err := accu.GetAccountList()
            if err != nil {
                http.Error(w, "email not found", http.StatusNotFound)               
                return
            }
            maxcount := len(accs)

            // loop through and return each field requested
            _, _ = fmt.Fprintf(w, "MaxCount: %d", maxcount)
            for _, acc := range accs {  //todo need some sensible sorting
                _,_ = fmt.Fprintf(w,"id: %s, email: %s",acc.ID,acc.Email)
                for i, field := range StringFields {
                    if _, found := r.Form[field]; found {
                        _, _ = fmt.Fprintf(w, ", %s: %s", field, u.StrUsecases[i])
                    } else {
                        // field not defined in the profile return a place holder
                        _, _ = fmt.Fprintf(w, ", %s: %s", field, "")
                    }
                }
                _,_ = fmt.Fprintf(w,"\n") //todo do as json instead of fprinting

                // todo images would go here either multipart/binary/urlencoded however its done in http
                for /*i*/ _, field := range ImageFields {
                    if _, found := r.Form[field]; found {
                        //_,_ = fmt.Fprintf(w,", %s: %s",field,u.ImageUsecases[i])
                    } else {
                        // todo no image defined case
                    }
                }
                _,_ = fmt.Fprintf(w,"\n") //todo do as json instead of fprinting
            }
            w.WriteHeader(http.StatusOK)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
        
    })
}
