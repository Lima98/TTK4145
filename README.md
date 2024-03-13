# Readme

## Important drivers and other files
https://github.com/TTK4145?q=driver

### Running the elevator
To run the simulator, navigate to the Elevator_project folder and run 
```
go run startSimulator.go
```

To run the physical elevator you need to specify "physical" mode and ID, the ID has to be unique on the network. To do this you navigate to Elevator_project and run
```
go run main.go physical "ID"
```
Where "ID" is replaced by a number. Convention here is to use 0,1,2 etc..