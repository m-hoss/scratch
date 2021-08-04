# scratch
a scratch discord bot for investigating gcp Cloud Run. 

since bot doesn't rely on incoming http traffic, currently there is an issue with more than one instance or revision being active. number of instances can be constrained in cloud run settings.

current issue is when new code is pushed and new build is deployed, existing revisions are not removed but reduced to 0% traffic. but since bot doesn't rely on traffic, multiple instances of bot will be running. currently only know how to solve this with manual revision deletion which is not ideal.

actually, maybe that's not the problem. 
