#!/usr/bin/env python
# coding: utf-8

# # Import desired utlities:
# 
# ### pip - conda installs if needed:
# pip install beautifulsoup4
# * User friendly and extendable html parser that makes it easy to find and loop through elements
# pip install lxml
# * beautiful soup uses external HTML/xml parsers so far this is the only one I've used
# https://www.crummy.com/software/BeautifulSoup/bs4/doc/
# 
# pip install requests
# * Used to fetch and get information about a web page request. One of the most userfriendly ways to get web page data
# 
# pip install dateparser
# * You will likely run into mutiple date time formats I've had good luck with this parser that is able to take a lot for formats of a string date and convert it to a python date time.
# 
# 
# 
# 

# In[29]:


import requests

from bs4 import BeautifulSoup
import dateparser

import json
import time


# In[2]:


#The business page you want to pull reviews from:
#Note* it seems like the business alias is the last part of the URL may be useful in scripting
BaseUrl = "https://www.yelp.com/biz/enamel-dentistry-south-lamar-austin"


# In[3]:


response = requests.get(BaseUrl)


# # Check For Valid webpage response
# * most webpage results will respond with status codes.
# * for most people 200 is good and means you got a good response everything else is "BAD" and an error
# * https://developer.mozilla.org/en-US/docs/Web/HTTP/Status
# * You will also notice that the text response on yelp is "nasty" to look at
# * If you do enough webscraping you can often tell what "framework" websites are using 
#     * in this case you see long lines, and generic and nonhuman readable parts of the CSS: list__09f24__ynIEd
#     * This usually means that a framework is being used, and potentially the CSS classes are generated dynamically with javascript or a backend language.

# In[4]:


if response.status_code == 200:
    print(response.text)


# In[5]:


#Parse the HTML content using lxml and put it into a BeautifulSoup object
SoupFullPage = BeautifulSoup(response.text, 'lxml')


# # Crash Course HTML parsing
# * Pretty much all webpages have "errors" and are not valid HTML if you were to trust a XML parser that relies on well opened closed, and structured XML.
# * id's **should** only apprear once per name per page. These are great ways to select exactly what you want
# * CSS classes are the next easiest thing. But these are usually not unique, so you can get odd results.
# * NOTE* a css class is defined without any spaces, if there are mutiple classes listed (spaces between classes) then you can use 1 or more classes to define what you are looking for:
#     * ul class="undefined list__09f24__ynIEd" there are 2 classes here a undefined one and one labeld list__09f24__ynIEd
# * If you can get close to the element you want the beautifulsoup find and find_all can return just a small part of the overall page that you can then again apply another find/find_all to. 
# * HTML is a hiearchy. you can navigate that hiearchy using parents or at the same level using siblings
#     * https://www.crummy.com/software/BeautifulSoup/bs4/doc/#find-parents-and-find-parent
# * If you can avoid it XPATH is one of the harder ways to navigate an HTML page and is not built into beautiful soup
#     * /html/body/yelp-react-root/div[1]/div[4]/div/div/div[2]/div/div[1]/main/div[3]/div/section/div[2]/div/div[5]
#     

# In[6]:


#Note* find returns the first occurance
#find_all returns a list of elements meeting the criteria even if only 1 is found
SoupIDreviews = SoupFullPage.find("div", {"id": "reviews"})
ListOfSoupULreviews = SoupIDreviews.find_all("ul", {"class": "undefined list__09f24__ynIEd"})


# In[7]:


SoupListReviews = ListOfSoupULreviews[0].find_all("li")


# In[8]:


#This is a good reason for find all as each person's review is a separate li element
for SoupReview in SoupListReviews:
    #If mutiple people have responded there can be more than one a tag, 
    #We assume that the first one is the person making the review
    Person = SoupReview.find("a", {"class": "css-1m051bw"})
    print("Review User ID:", Person['href'].split("userid=")[-1])
    print("  " + Person['href'])
    print("  " + Person.text.strip()) # strip removes any HTML tags
    
    #if you want more information about responses to review you might want mutiple dates
    ReviewDateList = SoupReview.find_all("span", {"class": "css-chan6m"})
    
    #Different ways you might want/need to deal with date/times
    #How a website handles timezones is something you might want to experiement with too.
    print("  " + ReviewDateList[0].text.strip())
    varReviewDate = dateparser.parse(ReviewDateList[0].text.strip())
    print("  ", varReviewDate, type(varReviewDate), varReviewDate.date(), type(varReviewDate.date()))
    
    #depending on the rating the other classes in this div can change so only 1 was used that *should* be in all of them
    SoupRating = SoupReview.find("div", {"class": "five-stars__09f24__mBKym"})
    print(SoupRating["aria-label"])
    print(SoupRating["aria-label"][0])
    
    #a lot of investigation went into this you can figure out the raving from the number of path results
    SVGPathList = SoupRating.find_all("path", {"opacity": "1"})
    
    print("    Expected Rating out of 10 from SVG: ", len(SVGPathList))
    for Path in SVGPathList:
        print("     Rating Color - " + Path['fill'])
        
    for ReviewText in SoupReview.find_all("p"):
        print("------------------------------")
        print(ReviewText.text.strip())
    print("===============================")
    


# In[9]:


#Cleanup the URL so it can be a file name
filename = response.url.replace("https://www.yelp.com/biz/", "")
filename = filename.replace("/", "_")
filename = filename.replace(":", " ")

filename = filename.replace("=", "")
filename = filename.replace("?", " ")


# If you are going to write the response out to a file ALWAYS use prettify it makes it possible to read it as a human.
with open(filename + ".html", "w") as f:
    f.write(SoupFullPage.prettify())


# In[22]:


SoupPagination = SoupFullPage.find("div", {"class": "pagination__09f24__VRjN4"})

print("Page navigation info")
print(SoupPagination.prettify())
print("-------------------")


#Get the first Div tag
NavigateSiblings = SoupPagination.div
print("nav Siblings (same hiearchy level)")
print(NavigateSiblings.prettify())
print("-------------------")

#Get the net sibling
PageInfoDiv = NavigateSiblings.find_next_sibling("div")

print(PageInfoDiv.prettify())
print("-------------------")

print(PageInfoDiv.text.strip())

CurrentPage_LastPage = PageInfoDiv.text.strip().split(" of ")


# In[23]:


print(CurrentPage_LastPage)


# In[47]:


def GetAndParseYelpReviewURL(url, saveHTML, TimesToTryAgain=3):
    #This will be the object this function returns
    ListDictResults = []
    
    #Try to go to the URL
    response = requests.get(url)
    
    #Check for a return status of 200 or repeat waiting 10 seconds between requests until try again is reached
    webpagefailures = 0
    while response.status_code != 200:
        time.sleep(10)
        response = requests.get(url)
        webpagefailures = webpagefailures + 1
        
        if webpagefailures > TimesToTryAgain:
            return "FAILURE"
    
    #Start parsing the html of the page:
    SoupFullPage = BeautifulSoup(response.text, 'lxml')
    
    # If you have selected to save the page then write it out to a file
    if saveHTML:
        #Cleanup the URL so it can be a file name
        filename = response.url.replace("https://www.yelp.com/biz/", "")
        filename = filename.replace("/", "_")
        filename = filename.replace(":", " ")

        filename = filename.replace("=", "")
        filename = filename.replace("?", " ")
        
        with open(filename + ".html", "w") as f:
            f.write(SoupFullPage.prettify())
    
    
    SoupIDreviews = SoupFullPage.find("div", {"id": "reviews"})
    ListOfSoupULreviews = SoupIDreviews.find_all("ul", {"class": "undefined list__09f24__ynIEd"})
    
    SoupListReviews = ListOfSoupULreviews[0].find_all("li")
    
    for SoupReview in SoupListReviews:
        ReviewDictionary = {}
        
        Person = SoupReview.find("a", {"class": "css-1m051bw"})
        
        ReviewDictionary['PersonID'] = Person['href'].split("userid=")[-1]
        ReviewDictionary['PersonURLpart'] = Person['href']
        ReviewDictionary['PersonName'] = Person.text.strip()
        
        #Only Getting the first date
        ReviewDateElement = SoupReview.find("span", {"class": "css-chan6m"})

        #If not storing in json better to use native python date
        #ReviewDictionary['ReviewDate'] = dateparser.parse(ReviewDateElement.text.strip())
        revdate = dateparser.parse(ReviewDateElement.text.strip())
        ReviewDictionary['ReviewDate'] = revdate.strftime('%Y-%m-%d')
                                   
        #rating                                                  
        SoupRating = SoupReview.find("div", {"class": "five-stars__09f24__mBKym"})
        
        ReviewDictionary['Rating'] = int(SoupRating["aria-label"][0])
        
        AllReviewText = ""
        for ReviewText in SoupReview.find_all("p"):
            print("------------------------------")
            AllReviewText = AllReviewText + " " + ReviewText.text.strip()
        ReviewDictionary['ReviewText'] = AllReviewText
                                                          
        ListDictResults.append(ReviewDictionary)

    return ListDictResults


# In[48]:


# For yelp in the address bar pages are repersented by the get variable &start=
# In our case page 1 does not have this page 2 is 10 page 3 is 20 etc

BusinessID = 12345
AllBusinessReviews = []
#Note this only records pages 2 and on to the dictionary
for i in range(2, int(CurrentPage_LastPage[1])+1):
    
    
    print("Page: ", i)
    print("start var: ", (i-1)*10)
    
    if "?" in BaseUrl:
        print("Already get variables just add to the end")
        NextURL = BaseUrl + "&start=" + str((i-1)*10)
    else:
        NextURL = BaseUrl + "?start=" + str((i-1)*10)
    
    #Get and Parse the next page of reviews
    CurrentPageReviews = GetAndParseYelpReviewURL(NextURL, saveHTML=True, TimesToTryAgain=3)
    AllBusinessReviews.extend(CurrentPageReviews)
    #for safty if the script fails export the current info to a json file every 10 pages
    if i%3 == 0:
        with open("ReviewBackup" + str(BusinessID) + ".json", "w") as outfile:
            json.dump({BusinessID: AllBusinessReviews, 'TotalReviews': len(AllBusinessReviews), 'PageNum': i}, outfile, indent=4)
        
with open(str(BusinessID) + ".json", "w") as outfile:
    json.dump({BusinessID: AllBusinessReviews, 'TotalReviews': len(AllBusinessReviews)}, outfile, indent=4)
        


# In[51]:


print(json.dumps({
    BusinessID: AllBusinessReviews, 
    'TotalReviews': len(AllBusinessReviews)
              }, indent=4))


# In[ ]:




