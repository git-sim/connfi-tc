package handlers

import (
    "fmt"
    "net/http"
    "github.com/git-sim/tc/app/usecase"
)

func HandleAccount(u usecase.AccountUsecase) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {  
        switch r.Method { 
        case http.MethodPost:
            email := r.URL.Query().Get("email")
            if email == "" {
                http.Error(w, "missing email name in query string", http.StatusBadRequest)
                return
            }
            err := u.RegisterAccount(email)
            if err != nil {                
                http.Error(w, "email not found", http.StatusBadRequest)             
                return
            }
            w.WriteHeader(http.StatusCreated)
        case http.MethodDelete:
            email := r.URL.Query().Get("email")
            if email == "" {
                http.Error(w, "missing email name in query string", http.StatusBadRequest)
                return
            }
            err := u.DeleteAccount(email)
            if err != nil {
                
                http.Error(w, "email not found", http.StatusNotFound)               
                return
            }
            w.WriteHeader(http.StatusOK)    
            
        case http.MethodGet:
            email := r.URL.Query().Get("email")
            if email == "" {
                http.Error(w, "missing email name in query string", http.StatusBadRequest)
                return
            }
            acc, err := u.GetAccount(email)
            if err != nil {
                http.Error(w, "email not found", http.StatusNotFound)               
                return
            }
            w.WriteHeader(http.StatusOK)
            _,_ = fmt.Fprintf(w,"id: %s, email: %s, FirstName: %s, LastName: %s\n",
                acc.ID,acc.Email,acc.FirstName,acc.LastName)
            //w.Write(acc)
                
        case http.MethodPut:
            email := r.URL.Query().Get("email")
            if email == "" {
                http.Error(w, "missing email name in query string", http.StatusBadRequest)
                return
            }
            http.Error(w, "Account edit Not impl", http.StatusNotImplemented)               
        
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
        
    })
}

func HandleAccountList(u usecase.AccountUsecase) http.Handler {
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
