#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Fri May 24 11:51:19 2019

@author: jeff
"""

import pandas
import matplotlib.pyplot as plt
import numpy
import math
from keras.models import Sequential 
from keras.layers import Dense
from keras.layers import Flatten
from keras.layers import LSTM
from sklearn.preprocessing import MinMaxScaler
from sklearn.metrics import mean_squared_error
from sklearn.model_selection import train_test_split
from sklearn.datasets.samples_generator import make_blobs


#columns to be used decided from basic classification in R
dataset = pandas.read_csv('sample.csv', usecols=[1,2,4,5,6,7,8,9,10,11,12,13,15,16,18,20], engine='python')
dataset_norm = (dataset - dataset.mean()) / (dataset.max() - dataset.min())
dataset_norm["class"] = dataset["class"]
#plt.plot(dataset_norm)

#set seed for reproducibility
numpy.random.seed(212919156)

#plt.figure()
#groups = [1,2,3,4,5,6,7,8,9,10,11,12,13,14,15]
#i = 1
#for group in groups:
#    plt.subplot(len(groups), 1, i)
#    plt.plot(dataset.values[:, group])
#    plt.title(dataset.columns[group], y=0.5, loc='right')
#    i+=1
#plt.show

def series_to_supervised(data, n_in=1, n_out=1, dropnan=True):
    n_vars =1 if type(data) is list else data.shape[1]
    df = pandas.DataFrame(data)
    cols, names = list(), list()
    #input sequence
    for i in range(n_in, 0, -1):
        cols.append(df.shift(i))
        names += [('var%d(t-%d)' % (j+1, i)) for j in range(n_vars)]
    #forecast sequence
    for i in range(0, n_out):
        cols.append(df.shift(-i))
        if i == 0:
            names += [('var%d(t)'%(j+1)) for j in range(n_vars)]
        else:
            names += [('var%d(t+%d)'%(j+1)) for j in range(n_vars)]
    #aggregate
    agg = pandas.concat(cols, axis=1)
    agg.columns = names
    #drop NaNs
    if dropnan:
        agg.dropna(inplace=True)
    return agg

scaler = MinMaxScaler(feature_range=(0,1))
scaled = scaler.fit_transform(dataset_norm.values.astype('float32'))
scaled[:,1] = scaled[:,1].astype('int')
reframed = series_to_supervised(scaled, 1, 1)
reframed.drop(reframed.columns[[16,18,19,20,21,22,23,24,25,26,27,28,29,30,31]], axis=1, inplace=True)

#split into train and test sets
train = reframed.values[1:-(math.floor(len(dataset)*0.3)),:]
test = reframed.values[-(math.floor(len(dataset)*0.3)):,:]
#split into inputs and outputs
train_X, train_y = train[:,:-1], train[:,-1]
test_X, test_y = test[:,:-1], test[:,-1]
#reshape for 3D input
train_X = train_X.reshape((train_X.shape[0], 1, train_X.shape[1]))
test_X = test_X.reshape((test_X.shape[0], 1, test_X.shape[1]))


#Recurrent Neural Network for sequence prediction
model = Sequential()
model.add(LSTM(100, input_shape=(train_X.shape[1], train_X.shape[2])))
model.add(Dense(1, activation='sigmoid'))
model.compile(loss='binary_crossentropy', optimizer='adam', metrics=['accuracy'])

history = model.fit(train_X, train_y, epochs=500, batch_size=128, validation_data=(test_X, test_y), verbose=2, shuffle=False)

plt.plot(history.history['loss'], label='train')
plt.plot(history.history['val_loss'], label='test')
plt.legend()
plt.show()

scores = model.evaluate(test_X, test_y, verbose=0)

# make a prediction
yhat = model.predict_classes(test_X)
truth = test_X[:,:,1].astype('int')

hits = 0;
misses = 0;

for i in range(len(truth)):
    if truth[i] == yhat[i]:
        hits = hits+1
    else:
        misses = misses+1
        
accuracy = hits/len(truth)

print('Measured Accuracy of Model: %.3f' % (accuracy*100))
print('Evaluated Accuracy of Model: %.3f' % (scores[1]*100) )



Xnew, _ = make_blobs(n_samples=1, centers=3, n_features=train_X.shape[2], random_state=1)
Xnew = Xnew.reshape((Xnew.shape[0],1,Xnew.shape[1]))
ynew = model.predict_classes(Xnew)

print("Predicting the Class of the next 1-event")
for i in range(len(Xnew)):
    print("Event %s, Predicted=%s" % (i, ynew[i]))



#use dummy data to predict the next few classes
Xnew, _ = make_blobs(n_samples=5, centers=3, n_features=train_X.shape[2], random_state=1)
Xnew = Xnew.reshape((Xnew.shape[0],1,Xnew.shape[1]))
ynew = model.predict_classes(Xnew)

print("Predicting the Class of the next 5-events")
for i in range(len(Xnew)):
    print("Event %s, Predicted=%s" % (i, ynew[i]))
    
Xnew, _ = make_blobs(n_samples=10, centers=3, n_features=train_X.shape[2], random_state=1)
Xnew = Xnew.reshape((Xnew.shape[0],1,Xnew.shape[1]))
ynew = model.predict_classes(Xnew)

print("Predicting the Class of the next 10-events")
for i in range(len(Xnew)):
    print("Event %s, Predicted=%s" % (i, ynew[i]))
