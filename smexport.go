package main

import (
    "net/http"
    "time"
    "encoding/json"
    "os"
    "fmt"
    
)


type ScopeLinePagedRecords struct {
    Num_results int  
    Page int
    Objects []ScopeLineRecord
    Total_pages int
}

type ScopeLineRecord struct {
    lineitem string
    id int
    project_id int
}



type UrlPagedRecords struct {
    Num_results int  
    Page int
    Objects []UrlRecord
    Total_pages int
}

type UrlRecord struct {
    Url string
    Id int
}


func main() {
    if os.Args[1] == "domain" {
        firstPage := new(ScopeLinePagedRecords) 
        link := "http://bbs-scopemanager-service:7000/api/scope_line"
        getJson("http://bbs-scopemanager-service:7000/api/scope_line?page=1", firstPage)
        totalPages := firstPage.Total_pages
        //fmt.Println(totalPages)
        
        for i := 1; i <= totalPages; i++ {
            //fmt.Println(i)
            concatenated := fmt.Sprintf("%s?page=%d", link, i)
            //fmt.Println(concatenated)
            
            jsonData := new(ScopeLinePagedRecords)
            getJson(concatenated, jsonData)
            
            for currentIndex := range jsonData.Objects {
                fmt.Println(jsonData.Objects[currentIndex].lineitem)
            }
        }
    }    
}
/*
func mainOLD() {
    if os.Args[1] == "domain" {
        firstPage := new(DomainPagedRecords) 
        link := "http://bbsstore-service:7002/api/dns_store"
        getJson("http://bbsstore-service:7002/api/dns_store?page=1", firstPage)
        totalPages := firstPage.Total_pages
        //fmt.Println(totalPages)
        
        for i := 1; i <= totalPages; i++ {
            //fmt.Println(i)
            concatenated := fmt.Sprintf("%s?page=%d", link, i)
            //fmt.Println(concatenated)
            
            jsonData := new(DomainPagedRecords)
            getJson(concatenated, jsonData)
            
            for currentIndex := range jsonData.Objects {
                fmt.Println(jsonData.Objects[currentIndex].Domain)
            }
        }
    }
    
    if os.Args[1] == "url" {
        firstPage := new(UrlPagedRecords) 
        link := "http://bbsstore-service:7002/api/url_store"
        getJson("http://bbsstore-service:7002/api/url_store?page=1", firstPage)
        totalPages := firstPage.Total_pages
        //fmt.Println(totalPages)
        
        for i := 1; i <= totalPages; i++ {
            //fmt.Println(i)
            concatenated := fmt.Sprintf("%s?page=%d", link, i)
            //fmt.Println(concatenated)
            
            jsonData := new(UrlPagedRecords)
            getJson(concatenated, jsonData)
            
            for currentIndex := range jsonData.Objects {
                fmt.Println(jsonData.Objects[currentIndex].Url)
            }
        }
    }
    
    if os.Args[1] == "param" {
        firstPage := new(UrlPagedRecords) 
        link := "http://bbsstore-service:7002/api/url_store"
        // we could improve this code by querying url_store api only for urls containing '?' and '='
        getJson("http://bbsstore-service:7002/api/url_store?page=1", firstPage)
        totalPages := firstPage.Total_pages
        //totalPages = 300
        
        //fmt.Println(totalPages)
        
        //var foundParams
        
        var foundParams []string
        
        for i := 1; i <= totalPages; i++ {
            //fmt.Println(i)
            concatenated := fmt.Sprintf("%s?page=%d", link, i)
            //fmt.Println(concatenated)
            
            jsonData := new(UrlPagedRecords)
            getJson(concatenated, jsonData)
            
            
            for currentIndex := range jsonData.Objects {
                currentObject  := jsonData.Objects[currentIndex]
                if strings.Contains(currentObject.Url, "?") {
                    //fmt.Println(currentObject.Url)
                    
                    paramStr := strings.Split(currentObject.Url, "?")[1]
                    params := strings.Split(paramStr, "&")
                    for _, param := range params {
                        // ignore any params without an =, these might be used for cache busting
                        if strings.Contains(param, "="){
                            paramKeyandValue := strings.Split(param, "=")
                            paramKey := paramKeyandValue[0]
                            //fmt.Println(paramKey)
                            if contains(foundParams, paramKey) == false {
                                fmt.Println(paramKey)
                                foundParams = append(foundParams, paramKey)
                            }
                        }
                    }
                    
                    
                }
            }
        }
    }
}
*/



func contains(arr []string, str string) bool {
   for _, a := range arr {
      if a == str {
         return true
      }
   }
   return false
}


var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
    r, err := myClient.Get(url)
    if err != nil {
        return err
    }
    defer r.Body.Close()

    return json.NewDecoder(r.Body).Decode(target)
}
