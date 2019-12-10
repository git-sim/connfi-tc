package handlers

import (
    "fmt"
    "net/http"
    "github.com/git-sim/tc/app/usecase"
    "strconv"
)

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

type pairFieldEnum struct {
    field string
    which int
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
        val := r.URL.Query().Get(k)
        if val != "" {
            numParsed++
            err := u.StrUsecases[i].Set(id64,val)
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
        _, ok := r.URL.Query()[k]
        if ok {
            val, err := u.StrUsecases[i].Get(id64)
            numParsed++
            if err != nil {
                http.Error(w, "err in get handleProfile string field",http.StatusBadRequest)
                numErrors++
                // report all errors
            }
            _,_ = fmt.Fprintf(w, "%s\n", val)
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
        email := r.URL.Query().Get("email")
        if email == "" {
            http.Error(w, "missing email name in query string", http.StatusBadRequest)
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
            numImagesParsed := 0
            numImagesErrors := 0
            //numImagesParsed, numImagesErrors := parseAndSetImageUsecase(id64,ImageFields,u,w,r)
            if (numStrsErrors + numImagesErrors == 0) && (numStrsParsed + numImagesParsed > 0)  {
                // we parsed at least 1 field, and there were not errors
                w.WriteHeader(http.StatusOK)
            }

        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
        
    })
}

func HandleProfileList(u usecase.AccountUsecase) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch r.Method { 
        case http.MethodGet:
            accs, err := u.GetAccountList()
            if err != nil {
                http.Error(w, "email not found", http.StatusNotFound)               
                return
            }
            w.WriteHeader(http.StatusOK)
            _,_ = fmt.Fprintf(w,"count: %d\n",len(accs))
            for _, acc := range accs {
                _,_ = fmt.Fprintf(w,"id: %s, email: %s\n",acc.ID,acc.Email)
            }
            //w.Write(acc)
                
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
        
    })
}
