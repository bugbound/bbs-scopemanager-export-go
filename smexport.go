package main

import (
    "net/http"
    "time"
    "encoding/json"
    "os"
    "fmt"
    "strings"
)

type TotalPageCountInfo struct {
    Total_pages int
}

type ScopeLinePagedRecords struct {
    Num_results int  
    Page int
    Objects []ScopeLineRecord
    Total_pages int
}

type ScopeLineRecord struct {
    Id int
    Lineitem string
    Project_id int
}


type DnsStorePagedRecords struct {
    Num_results int  
    Page int
    Objects []DnsStoreRecord
    Total_pages int
}

type DnsStoreRecord struct {
    Id int
    Domain string
}

type IpStorePagedRecords struct {
    Num_results int  
    Page int
    Objects []IpStoreRecord
    Total_pages int
}

type IpStoreRecord struct {
    Id int
    Domain string
    Ip string
}


func main() {
    if os.Args[1] == "domain" {
        var scope_lines []string
        firstPage := new(ScopeLinePagedRecords) 
        
        BbsProjectId := os.Args[2]
        link := "http://bbs-scopemanager-service:7000/api/scope_line"
        searchQuery := fmt.Sprintf(`{"filters":[{"name":"project_id","op":"eq","val":"%s"}]}`, BbsProjectId)
        firstPageLink := fmt.Sprintf(`%s?page=1&q=%s`, link, searchQuery)
        
        // get page count
        getJson(firstPageLink, firstPage)
        totalPages := firstPage.Total_pages
        //fmt.Println(totalPages)
        
        for i := 1; i <= totalPages; i++ {
            concatenated := fmt.Sprintf("%s?page=%d&q=%s", link, i, searchQuery)
            
            jsonData := new(ScopeLinePagedRecords)
            getJson(concatenated, jsonData)
            
            //fmt.Println("Getting line items...")
            for currentIndex := range jsonData.Objects {
                if contains(scope_lines, jsonData.Objects[currentIndex].Lineitem) == false {
                    scope_lines = append(scope_lines, jsonData.Objects[currentIndex].Lineitem)
                }
            }
        }
            
        // now loop scope getting domains
        for currentIndex := range scope_lines {
            var Domains = getDomainListFromWildcardScopeLine(scope_lines[currentIndex])
            fmt.Println(strings.Join(Domains, "\n"))        
        }
    }    


    if os.Args[1] == "externalip" {
        var scope_lines []string
        var allDomains []string
        var allIps []string
        firstPage := new(ScopeLinePagedRecords) 
        link := "http://bbs-scopemanager-service:7000/api/scope_line"
        getJson(link, firstPage)
        totalPages := firstPage.Total_pages
        fmt.Println(totalPages)
        
        for i := 1; i <= totalPages; i++ {
            //fmt.Println(i)
            concatenated := fmt.Sprintf("%s?page=%d", link, i)
            //fmt.Println(concatenated)
            
            jsonData := new(ScopeLinePagedRecords)
            getJson(concatenated, jsonData)
            
            //fmt.Println("Getting line items...")
            for currentIndex := range jsonData.Objects {
                if contains(scope_lines, jsonData.Objects[currentIndex].Lineitem) == false {
                    scope_lines = append(scope_lines, jsonData.Objects[currentIndex].Lineitem)
                }
            }
        }
            
        // now loop scope getting domains
        for currentIndex := range scope_lines {
            //fmt.Println(scope_lines[currentIndex])
            var Domains = getDomainListFromWildcardScopeLine(scope_lines[currentIndex])
            for currentDomainIndex := range Domains {
                if contains(allDomains, Domains[currentDomainIndex]) == false {
                    allDomains = append(allDomains, Domains[currentDomainIndex])
                    var ipRecords = getIPListFromDomain(Domains[currentDomainIndex])
                    
                    for iprIndex := range ipRecords {
                        isExternal := true
                        if strings.HasPrefix(ipRecords[iprIndex], "192.168") == true {isExternal=false}
                        if strings.HasPrefix(ipRecords[iprIndex], "10.") == true {isExternal=false}
                        
                        if isExternal == true {
                            if contains(allIps, ipRecords[iprIndex]) == false {
                                if isIpOnIgnoreList(ipRecords[iprIndex]) == false {
                                    fmt.Println(ipRecords[iprIndex])
                                    allIps = append(allIps, ipRecords[iprIndex])
                                }
                            }
                        }
                    }
                }
            }
        }        
        
    }    
    
    if os.Args[1] == "internalip" {
        var scope_lines []string
        var allDomains []string
        var allIps []string
        firstPage := new(ScopeLinePagedRecords) 
        link := "http://bbs-scopemanager-service:7000/api/scope_line"
        getJson(link, firstPage)
        totalPages := firstPage.Total_pages
        fmt.Println(totalPages)
        
        for i := 1; i <= totalPages; i++ {
            //fmt.Println(i)
            concatenated := fmt.Sprintf("%s?page=%d", link, i)
            //fmt.Println(concatenated)
            
            jsonData := new(ScopeLinePagedRecords)
            getJson(concatenated, jsonData)
            
            //fmt.Println("Getting line items...")
            for currentIndex := range jsonData.Objects {
                if contains(scope_lines, jsonData.Objects[currentIndex].Lineitem) == false {
                    scope_lines = append(scope_lines, jsonData.Objects[currentIndex].Lineitem)
                }
            }
        }
            
        // now loop scope getting domains
        for currentIndex := range scope_lines {
            //fmt.Println(scope_lines[currentIndex])
            var Domains = getDomainListFromWildcardScopeLine(scope_lines[currentIndex])
            for currentDomainIndex := range Domains {
                if contains(allDomains, Domains[currentDomainIndex]) == false {
                    allDomains = append(allDomains, Domains[currentDomainIndex])
                    var ipRecords = getIPListFromDomain(Domains[currentDomainIndex])
                    
                    for iprIndex := range ipRecords {
                        isExternal := true
                        if strings.HasPrefix(ipRecords[iprIndex], "192.168") == true {isExternal=false}
                        if strings.HasPrefix(ipRecords[iprIndex], "10.") == true {isExternal=false}
                        
                        if isExternal == false {
                            if contains(allIps, ipRecords[iprIndex]) == false {
                                fmt.Println(ipRecords[iprIndex])
                                allIps = append(allIps, ipRecords[iprIndex])
                            }
                        }
                    }
                }
            }
        }
        
        //for currentIpIndex := range allIps {
        //    fmt.Println(allIps[currentIpIndex])
        //}
        
        //fmt.Println("Getting hosts for line items")
        
        
        
    }   

}

type IpOnIgnoreListRecord struct {
    Ip string
    Ignore bool
}

func isIpOnIgnoreList(ipToCheck string) bool  {
    link := fmt.Sprintf("http://192.168.26.1:7000/check_ip?ip=%s", ipToCheck)
    result := new(IpOnIgnoreListRecord) 
    getJson(link, result)
    return result.Ignore
}

func getDomainListFromWildcardScopeLine(scopeline string) []string {
    
    var domain string = strings.Replace(scopeline, "*", "%", -1)
    var records []string
    records = append(records, strings.Replace(scopeline, "*.", "", -1))
    //fmt.Println(domain)
    link := "http://bbsstore-service:7002/api/dns_store"
    searchQuery := fmt.Sprintf(`{"filters":[{"name":"domain","op":"like","val":"%s"}]}`, domain)
    firstPageLink := fmt.Sprintf(`%s?page=1&q=%s`, link, searchQuery)
    //fmt.Println(firstPageLink)

    totalPageCount := new(TotalPageCountInfo) 
    getJson(firstPageLink, totalPageCount)

    totalPages := totalPageCount.Total_pages
    //fmt.Println(totalPages)
    
    for i := 1; i <= totalPages; i++ {
            concatenated := fmt.Sprintf("%s?page=%d&q=%s", link, i, searchQuery)
            //fmt.Println(concatenated)
            jsonData := new(DnsStorePagedRecords)
            getJson(concatenated, jsonData)     
            for currentIndex := range jsonData.Objects {
                lookup := jsonData.Objects[currentIndex].Domain
                if strings.Index(lookup, "*.") > -1 {
                    lookup = strings.Replace(lookup, "*.", "", -1)
                }
                
                if contains(records, lookup) == false {
                    records = append(records, lookup)
                }
            }        
        }
    
    return records
}


func getIPListFromDomain(scopeline string) []string {
    
    var domain string = strings.Replace(scopeline, "*", "%", -1)
    var records []string
    //fmt.Println(domain)
    link := "http://bbsstore-service:7002/api/ip_store"
    searchQuery := fmt.Sprintf(`{"filters":[{"name":"domain","op":"eq","val":"%s"}]}`, domain)
    firstPageLink := fmt.Sprintf(`%s?page=1&q=%s`, link, searchQuery)
    //fmt.Println(firstPageLink)

    totalPageCount := new(TotalPageCountInfo) 
    getJson(firstPageLink, totalPageCount)

    totalPages := totalPageCount.Total_pages
    //fmt.Println(totalPages)
    
    for i := 1; i <= totalPages; i++ {
            concatenated := fmt.Sprintf("%s?page=%d&q=%s", link, i, searchQuery)
            //fmt.Println(concatenated)
            jsonData := new(IpStorePagedRecords)
            getJson(concatenated, jsonData)     
            for currentIndex := range jsonData.Objects {
                lookup := jsonData.Objects[currentIndex].Ip
                
                if contains(records, lookup) == false {
                    records = append(records, lookup)
                }
            }        
        }
    
    return records
}



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
    //fmt.Println("response Status:", r.Status)
    
    return json.NewDecoder(r.Body).Decode(target)
}





