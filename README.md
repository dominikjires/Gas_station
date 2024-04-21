# Gas_station
Go concurrency  
Create gas station simulation  

1. Cars arrive at the gas station and wait in the queue for the free station
2. Total number of cars and their arrival time is configurable
3. There are 4 types of stations: gas, diesel, LPG, electric
4. Count of stations and their serve time is configurable as interval (e.g. 2-5s) and can be different for each type
5. Each station can serve only one car at a time, serving time is chosen randomly from station's interval
6. After the car is served, it goes to the cash register.
7. Count of cash registers and their handle time is configurable
8. After the car is handled (random time from register handle time range) by cast register, it leaves the station.
9. Program collects statistics about the time spent in the queue, time spent at the station and time spent at the cash register for every car
10. Program prints the aggregate statistics at the end of the simulation

Spuštění docker image: ```docker run jiresdom/gasstation```  
Repozitář: https://hub.docker.com/r/jiresdom/gasstation
