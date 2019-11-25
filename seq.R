library(glmnet)
library(neuralnet)
library(NeuralNetTools)
library(nnet)
library(randomForest)

set.seed(212919156)

#read in the sample data
sample <- read.csv("seq/sample.csv")
sample$class = factor(sample$class)

#Use a random forest to generate a model to predict class
rf <- randomForest(factor(class)~., data=data.frame(sample[,-1]))

#Now look at the Variable Importance Plot of the Random Forest to determine relevant variables to the model
varImpPlot(rf)

#based on the variable importance plot, the following variables are below an arbitrary 200 MeanGiniDecrease threshold
#So, we'll disregard them and remake the data without 
simdat <- sample[,-c(1,4,24,17,21,23,14,19,22,25,26,27,28,29,30,31)]

#now using only those, we can run a logistic regression
simlog <- glm(class ~ ., family="binomial", data=simdat)

#Now get the confusion matrix for both
#First the randomForest
table(predict(rf, newdata=sample[,-1]), sample$class)

#Then the logistic regression
table(predict(simlog, type="response") > 0.5, sample$class)

#obviously, the randomForest performs better for predictions, so given data about the next events, we would use the randomForest
#model to predict the class. Random Forests have built in cross-validation as well, due to Out-of-Bag estimation, 
#so there's no need to additionally cross-validate this model. 

