# %%
# Example of Pulling Json API data. 12-19-2022
import cred_user_pwd_keys # a .py file in the same folder that has account details
# Other ways to deal with credentials.
#  - If you trust the computer and are performing development to just make sure things work.
#     * the easisest approach is to embed the API keys/passwords/usernames in this file.
#     * I do not do this, but I do create a copy of the credential information into a specific variable
#     * you can just replace that definition with your credentials
#     * if you store your passwords in plain text in this file DO NOT unpload to git hub or paste the code anywhere on the internet
# - The docker approache: Note the user credentials can be stored as enviroment variables (often used in docker setups for better security)
#     * https://www.twilio.com/blog/environment-variables-python
#     * https://chlee.co/how-to-setup-environment-variables-for-windows-mac-and-linux/
#     * https://blog.pilosus.org/posts/2019/06/07/application-configs-files-or-environment-variables-actually-both/
# - The final way and the way I'm handling this is creating a "settings/config/passwords" text file that stores the information in plain text
#     * THIS settings/passwords file should NEVER be uploaded to Git
#     * add this file to your .gitignore file or similar 
#DO NOT UPLOAD this file if you have hardcoded API keys/passwords
#DO NOT UPLOAD THE cred_user_pwd_keys.py file to GIT or any other file sharing site
import os
import shutil # for copying .jpgs or other raw data from the web request
# Requests is one of the easier to understand modules to download a website or many web related APIs
import requests


#I don't put much into pandas, but you can load info into pandas for all your data manipulation needs
import pandas as pd
# Although not scrictly needed most API calls will return a json format. 
# Python's dictionary and Json are very simlar and easy to swtich between.
import json
import re # used for pattern matching and removing bad characters


#Dealing with dates/times
from datetime import datetime
from dateutil.parser import parse #Let the parser deal with the formatting
import time 

#Example of how to pull enviroment variables
PythonLocation = os.environ.get('PATH')

if PythonLocation:
    print('All your application locations: ')
    print(PythonLocation)
else:
    print("you don't have an enviroment PATH")

# %% [markdown]
# # EASY KEY IS PART OF ADDRESS
# ## The Movie DB V3
# * Always refer back to the API information.
# * More information can be found here
# 
# https://developers.themoviedb.org/3/getting-started/introduction

# %%
#The Movie DB V3 Example pulled right from the API information
movieID = 550
#Manual Test Example
akey = cred_user_pwd_keys.TheMovieDB_API_Key_v3_auth

MainMovieInfoUrl = f'https://api.themoviedb.org/3/movie/{movieID}?api_key={akey}'


testWebGet = requests.get(MainMovieInfoUrl)

ParsedData = json.loads(testWebGet.text)

#View the result in a more human readable format
print(json.dumps(ParsedData, indent=4))

# %%
print("Title: ", ParsedData['title'], "| SecondProdCompany: " + ParsedData['production_countries'][1]['name'])
print("Looping through all Genre's")
for genre in ParsedData['genres']:
    print("   sub-json data:", genre)
    print("   -  ", genre["name"])

# %%
#List all production Companies:
for comp in ParsedData['production_companies']:
    print(comp['name'], ' From: ', comp['origin_country'])

# %%
#Export Results to a .json file as an example of one way to extract information 
with open(f"{ParsedData['imdb_id']}.json", 'w') as fp:
    json.dump(ParsedData, fp)
    print("Saved File to: " + os.getcwd() + f"/{ParsedData['imdb_id']}.json")

# %%
# Using the Find with a IMDb ID---


movie_IMDb_id = 'tt4154796'
externalsource= 'imdb_id'


url = f'https://api.themoviedb.org/3/find/{movie_IMDb_id}?api_key={akey}&language=en-US&external_source={externalsource}'

webdata = requests.get(url)
ParsedData = json.loads(webdata.text)

#View the result in a more human readable format
print(json.dumps(ParsedData, indent=4))

GetTheID = ParsedData["movie_results"][0]['id']

print("TheFirstMovieFoundID:", GetTheID)

# %%
# Get one reviews

movie_id = '299534'

pageLoop = 1

url = f'https://api.themoviedb.org/3/movie/{movie_id}/reviews?api_key={akey}&language=en-US&page={pageLoop}'

webdata = requests.get(url)
ParsedData = json.loads(webdata.text)

#View the result in a more human readable format
print(json.dumps(ParsedData, indent=4))



# %%
print("Number of Separate Web Calls (pages) Needed: " + str(ParsedData['total_pages']))
print(" - Current (page): " + str(ParsedData['page']))
print(" - Number of descrete reivews on this (page): " + str(len(ParsedData['results'])))

print("")

CreationDTobject = parse(ParsedData['results'][0]['created_at'])
print("The Date time is a: " + str(type(CreationDTobject)))

print("The First Review was created on: " + str(CreationDTobject.date()))
print("    At this time: " + str(CreationDTobject.time()))
print("    In the TimeZone: " + str(CreationDTobject.tzinfo))
print("    or human readable: " + CreationDTobject.strftime("%A %m/%d/%Y, %I:%M %p %Z"))
#https://www.programiz.com/python-programming/datetime/strftime

# %% [markdown]
# ## Queries with paramaters
# 

# %%
#Example pulled from the API's examples
pageLoop = 1

#We want all movies that are less than PG-13 rating us rating.
#To get more realistic results only return movies with more than 75 votes
#Sort by the vote_average descending. This will give us A LOT of movies it would be worth adding more filters
certification_country="US"
certificationMaxScore="PG-13"
GTNumbVotes=75
sort_by="vote_average.desc"


url = f'https://api.themoviedb.org/3/discover/movie/?api_key={akey}&language=en-US&page={pageLoop}&certification_country={certification_country}&ertification.lte={certificationMaxScore}&vote_count.gte={GTNumbVotes}&sort_by={sort_by}'

webdata = requests.get(url)
ParsedData = json.loads(webdata.text)

#View the result in a more human readable format
print(json.dumps(ParsedData, indent=4))



# %%
#Take the example of the query above and pull the 3 & 4th page of results only getting the movie ID so we can extract the reviews of those movies

ListOfMovieIDs=[]

#Grab all the movie Ids on the 3 & 4th page
#Example pulled from the API's examples
for pagenum in range (2, 4):
    pageLoop = pagenum

    certification_country="US"
    certificationMaxScore="PG-13"
    GTNumbVotes=75
    sort_by="vote_average.desc"


    url = f'https://api.themoviedb.org/3/discover/movie/?api_key={akey}&language=en-US&page={pageLoop}&certification_country={certification_country}&ertification.lte={certificationMaxScore}&vote_count.gte={GTNumbVotes}&sort_by={sort_by}'

    webdata = requests.get(url)
    ParsedData = json.loads(webdata.text)
    #Loop through each movie
    for Movie in ParsedData['results']:
        ListOfMovieIDs.append(Movie['id'])
        #Note These print statments might fail as there are non US characters in a lot of the movie names
        print("movie: ", Movie['original_title'], " Vote AVG: ", str(Movie['vote_average']), " Vote Count: ", str(Movie['vote_count']))




# %% [markdown]
# ### an example of taking the input movie IDs from above's list and an approach to save this to a file

# %%
#Setup some variables
MoviewReviewsList = []
output = pd.DataFrame()
movieCount = 0


for movieid in ListOfMovieIDs:
    movieCount = movieCount + 1
    
    #Do something like this if the run failed and you need to start somewhere in the middle. This will waste some time looping through movies, but it works
    #if movieCount <= 2600:
    #    continue
   
    #Used to occassionally save the output. If there is a crash we can just start about where it failed
    pageLoop = 1
    reviews = []
    
    # We have to call the API once to see how many pages of reviews there are
    url = f'https://api.themoviedb.org/3/movie/{movieid}/reviews?api_key={akey}&page={pageLoop}'
    webdata = requests.get(url)
    ParsedData = json.loads(webdata.text)

    # If there are no results or bad movie go to the next movie
    if ParsedData.get('results') == None:
        #Did not get a response for the movie
        continue
    #Remember the structure above the results is a list of dictionary reviews, we append this to the temporary python list reviews
    reviews.append(ParsedData['results'])
    
    # Now that we have the initial pull if there is more than 1 page of reviews loop through all pages. Note for Large reviews in the 1000s you might want
    # To cut it off early like check if total pages > 100 and stop it there.
    if ParsedData['total_pages'] > 1:
        totalpages = ParsedData['total_pages']
        print("MoreThan1page ", movieid)
        for i in range(2, totalpages+1):
            pageLoop = i
            webdata = requests.get(url)
            ParsedData = json.loads(webdata.text)
            reviews.append(ParsedData['results'])
    MoviewReviewsList.append({"id":movieid, "Reviews": reviews})
    
    #Every 25 movies print off the status
    if movieCount % 25 == 0:
        print(movieid, ' On Movie: ', str(movieCount))
    
    #Every 30 movies Save a backup of the currently completed pulls.
    #Very useful to not keep pulling the same info, but there is still a chance you will get duplicate records that you will need to clean.
    if movieCount % 30 == 0:
        with open(f"5000_TMDB_Reviews{movieCount}.json", 'w') as fp:
            json.dump(MoviewReviewsList, fp)

# %%
#Assuming the above code completed without error... Save a final copy of the data

with open(f"001_ALL_TMDB_Reviews.json", 'w') as fp:
    json.dump(MoviewReviewsList, fp)

# %% [markdown]
# ## Another Example of Paramaters and saving the data as it goes

# %%
#Only english movies, with original language as english
#Sort by Newest to Oldest primary release date
#Only get movies with over 200 votes
#DOES NOT USE THE POPULARITY NUMBERS
#Because of this is the discover API we get slightly less info on the movies.
lang = 'en-US'
sort = 'primary_release_date.desc'
votecnt = '200'
origlang = 'en'

#Lets store all the info in a dataframe.
output = pd.DataFrame()

movieCount = 0
for page in range(1, 50):
    movieCount = movieCount + 1

    
    #This API is a generic search API to returna list of movies that meet the criteria
    movieurl = f'https://api.themoviedb.org/3/discover/movie?api_key={cred_user_pwd_keys.TheMovieDB_API_Key_v3_auth}&language={lang}&sort_by={sort}&include_adult=false&include_video=false&page={page}&vote_count.gte={votecnt}&with_original_language={origlang}'

    #Go to the page and load it into json
    moviepage = requests.get(movieurl)
    ParsedData = json.loads(moviepage.text)
    
    #Just incase we don't get results wait and re-try
    if ParsedData.get('results') == None:
        #Did not get a response Wait 120 seconds and continue to next page
        time.sleep(120)
        continue
    for movie in ParsedData['results']:
        #Keep dumping retulsts list of movies with some info into the output dataframe
        #TODO look into a better way to do this...
        output = output.append(movie, ignore_index=True)
    
    #Print status to screen every 10 pages
    if movieCount % 10 == 0:
        print(ParsedData['results'][0]['original_title'], ' On Page: ', str(movieCount))
    #Save a running file every 20 pages for safey incase the code crashes
    if movieCount % 20 == 0:
        output.to_csv(f"RecentPopularMovies{movieCount}.csv")
        output.to_pickle(f"RecentPopularMovies{movieCount}.pkl")     
    
    if movieCount % 25 == 0:
        time.sleep(60) #Wait 1 minutes every 25 pages to not overload API this seems like overkill could probably increase to 100
    
    
output.to_csv(f"CompleteRecentPopularMovies.csv")
output.to_pickle(f"CompleteRecentPopularMovies.pkl")

# %% [markdown]
# ## Data outside API like Images or additional links
# 

# %%
#Download an Image from the API many APIs will have arbitrary ways to pull extra information like images or links to additional information.
# If you are lucky this will be defined in the API documentation, or in an API call itself as is the cass with the movie db
#https://developers.themoviedb.org/3/getting-started/images
    
ManuallyDefinedUrl = f'https://api.themoviedb.org/3/configuration?api_key={cred_user_pwd_keys.TheMovieDB_API_Key_v3_auth}'

testWebGet = requests.get(ManuallyDefinedUrl)

ParsedConfig = json.loads(testWebGet.text)

#View the result in a more human readable format
print(json.dumps(ParsedConfig, indent=4))

# %%
# Get one reviews

movie_id = '299534'

pageLoop = 1

url = f'https://api.themoviedb.org/3/movie/{movie_id}?api_key={akey}&language=en-US'

webdata = requests.get(url)
ParsedData = json.loads(webdata.text)





# %%
print(ParsedData)

# %%
#View the result in a more human readable format
for key, value in ParsedData.items():
    print('Key:', key, ' | Value: ', value)
    print("--------")

# %%
# Download Image with info from the API and info in the results
baseURL = 'http://image.tmdb.org/t/p/'
poster_size = 'w342' #  Other Options "w92", "w154", "w185", "w342", "w500"

FileName = ParsedData['poster_path']

Savepath = "MoviePosterFromWeb.jpg"


poster_full_path_url = baseURL + poster_size + FileName

print("Full Url for Image: ", poster_full_path_url)

r = requests.get(poster_full_path_url, stream=True)
if r.status_code == 200: # in web speak 200 means a good response that we can read
    with open(Savepath, 'wb') as f:
        r.raw.decode_content = True
        shutil.copyfileobj(r.raw, f)     
        print('Image sucessfully Downloaded: ',Savepath) 

# %%



