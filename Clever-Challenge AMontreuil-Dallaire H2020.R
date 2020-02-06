#Part 3
#The third part of the challenge investigates sequences. In the seq folder, you will find two csv files. 
#The first file sample.csv contains ~18k events that are classified into two classes: 0 and 1. 
#All events are uniquely identified by their id and occurs at a precise timestamp. 
#In addition to the id, the timestamp and the class each event is further categories using 30 metrics (f1, ..., f30).

#Another file, named res.csv contains ids of the ressources used by the events describeb in sample.csv.

#In this part, we are interested in predicting, in advance, what will be the class of the next 1-,5- and 10- events, 
#given the information known at a precise time. For example, using the all relevent information we know after event 829454 
#(f1 to f30 and ressources information) can you predict the class of event 829455, 829455@829459 and 829455@829464.

#You are free to use the technics / languages of your choice.



#Set-up the data

require(Matrix)
require(ISLR)

#Sample data
sampleData <- read.csv(file='seq/sample.csv')
sData <- sampleData
rownames(sData) <- paste('ID', sData[,1], sep='')

#Ressources data
resData <- read.csv(file='seq/res.csv')
resSM <- sparseMatrix(resData[,1],resData[,2])

#Reshape the sparse matrix to the scope of ids and ressources we have
rownames(resSM) <- paste('ID', 1:dim(resSM)[1], sep='')
colnames(resSM) <- paste('RES', 1:dim(resSM)[2], sep='')
res <-resSM[rowSums(resSM)>0,colSums(resSM)>0]


##Functions

#Determines the similarities between the events with the cosinus formula
#v : vector of the event we want to predict
#m : matrix of all the other events that has at least 1 common ressource with the event we want to predict
#thr : threshold to determine if an event is similar or not
#return : a vector of names of the event that are at least similar as the threshold
cosinusVM <- function(v, m, thr) {
  n <- sqrt(colSums(m^2))
  cos <- (v %*% m)/(n * sqrt(sum(v^2)))
  cos <- as.matrix(cos)
  colnames(cos)<-colnames(m)
  simi <- colnames(cos)[cos>thr]
  return(simi)
}


#Get the similar events based on the common ressources with the given event, evaluated with the cosinus function.
#res : sparse matrix of the ressources data
#id : name of the event
#thr : threshold to determine if an event is similar or not
#return : a vector of names of the most similar events
CommonRessources <- function(res, id, thr){
  resEvent <- res[id,]>0
  resTF <- t(res>0)
  resIdEvent <- t(resTF*resEvent)
  similarIds <- res[rowSums(resIdEvent)>0,]
  common <- similarIds[,colSums(similarIds)>0]
  if(!(class(common) == "ngCMatrix")){
    common <- as.matrix(common)
    colnames(common) <- colnames(similarIds)[colSums(similarIds)>0]
  }
  commonT <- t(common)
  ordo <- cosinusVM(commonT[,id], commonT, thr)
  return(ordo)
}

#Verifies if the predicted class correspond to the real class of the event
#sData : Samples data with all the features of the events
#id : name of the event
#classPredict : predicted class
#Return : boolean saying if we predicted right
ResultPrediction <- function(sData, id, classPredict){
  classReal <- sData[id, 3]
  return(classReal == classPredict)
}

#Estimates the class of a given event, based on a trained logistic regression model on the similar events
#event : number of the event we want to predict
#sData : Samples data with all the features of the events
#res : sparse matrix of the ressources data
#thr : threshold to determine if an event is similar or not
EstimateClassResult <- function(event, sData, res, thr){
  id <- paste('ID', event, sep='')
  # If the event has no ressources in common with any other events, return 0
  if(sum(res[id,])){
    if(sum(res[,colnames(res)[res[id,]>0]])==sum(res[id,]>0)){
      return(0)
    }
  }
  similarID <- CommonRessources(res, id, thr)
  #If the only id enough similar is itself, we can't predict the class
  if (sum(res[similarID,]) == sum(res[id,])){
    return(0)
  }
  sampleTrain <- sData[similarID,]
  sampleTrain <- sampleTrain[-which(rownames(sampleTrain) == id),]
  sampleTest <- sData[id,]
  sampleLogistic <- glm(class ~ f1 + f2 + f3 + f3.1 + f4 + f5 + f6 + f7 + f8 + f9 + f10 + f11 + f12 + f13 + f14 + f15 + f16 + f17 + f18 + f19 + f20 + f21 + f21, data = sampleTrain, family = binomial)
  idProb <- predict(sampleLogistic, newdata = sampleTest, type = "response")
  classPredict <- (idProb>0.5)*1
  return(classPredict)
}


#Classifies an array of event
#sData : Samples data with all the features of the events
#res : sparse matrix of the ressources data
#thr : threshold to determine if an event is similar or not
#Return : a vector of the predicted classes
Classificator <- function(sData, res, thr, part){
  classPredict <- rep(0, length(part))
  i <- 1
  for (p in part){
    id <- paste('ID', p, sep='')
    if (!is.element(id, rownames(res))){
      classPredict[i] <- 0
    } else {
      classPredict[i] = EstimateClassResult(p, sData, res, thr)
    }
    i <- i + 1
  }
  names(classPredict) <- paste('ID', part, sep='')
  return(classPredict)
}


#Calculates the results in percentage depending on the threshold of similar event
#thrs : vector of thresholds to test
#sData : Samples data with all the features of the events
#res : sparse matrix of the ressources data
#part : vector of events randomly sampled
#return : the results of precisions for each tested threshold
OptimizeThreshold <- function(thrs, sData, res, part){
  results <- rep(0, length(thrs))
  i <- 1
  for (thr in thrs){
    classesPredict <- Classificator(sData, res, thr, part)
    result <- 0
    j <- 1
    for (id in names(classesPredict)){
      result <- result + unname(ResultPrediction(sData, id, classesPredict[id]))
      j <- j + 1
    }
    results[i] = result/length(part)
    i <- i + 1
  }
  return(results)
}

#Here is how I determined the optimal threshold

#thrs <- c(0.05, 0.1, 0.2, 0.35, 0.5)
#partThr <- sample(sData[,"event_id"], 2000)
#thrOpt <- OptimizeThreshold(thrs, sData, res, partThr)
#thrOpt
#optResults <- cbind(thrs, thrOpt)
#plot(optResults)

# I choose thr = 0.35
thr <- 0.35


#Predictes the class of the given event and the next ones
#event : number of the first event we want to predict
#sData : Samples data with all the features of the events
#res : sparse matrix of the ressources data
#thr : threshold to determine if an event is similar or not
#f : number of future events we want to predict
#Return : None. The function prints results on the console
PredictingFutureEvents <- function(event, sData, res, thr, f){
  cat("Predictions for the next events :")
  cat("\n")
  i <- 1
  ev <- event
  predict <- rep(0, f)
  class <- rep(0, f)
  while (i<(f+1)){
    id <- paste('ID', ev, sep='')
    if (!is.element(id, rownames(res))){
      predict[i] <- 0
    } else {
      predict[i] <- EstimateClassResult(ev, sData, res, thr)*1
    }
    names(predict)[i] <- id
    class[i] <- c(sData[id, 3])
    names(class)[i] <- id
    ev <- ev + 1
    i <- i + 1
  }
  print(predict)
  cat("Real class of the events :")
  cat("\n")
  print(class)
}

options(warn = -1)

#################################################
# RESULTS
#################################################


# Performance of the algorithm (based on a random sample or on the total events)
total = dim(sData)[1]
nbr = 2000
part <- sample(sData[,"event_id"], nbr)

classPredict <- Classificator(sData, res, thr, part)
classPredict

classTable <- table(classPredict, sData[paste('ID', part, sep=''), "class"])
classTable

recall <- classTable[1,1]/sum(classTable[,1])
recall
precision <- classTable[1,1]/sum(classTable[1,])
precision


# Recall : 61,0% (based on 2000 random events)
# Precision : 69,2% (based on 2000 random events)



#Write the event here for future event predictions
event <- 811355

#Run these function to see the next predictions
PredictingFutureEvents(event, sData, res, thr, 1)
PredictingFutureEvents(event, sData, res, thr, 5)
PredictingFutureEvents(event, sData, res, thr, 10)


