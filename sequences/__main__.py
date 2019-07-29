import numpy
import numpy as np
import pandas as pd
from keras.layers import LSTM, Dense
from keras.models import Sequential
from sklearn.metrics import mean_squared_error
from matplotlib import pyplot

def predict_sequence():
    sample = pd.read_csv("../seq/sample.csv")
    # ressources = pd.read_csv("../seq/res.csv")
    sample = prepare_data(sample)

    # split source data in training and validation sets
    x_train, y_train, x_validate, y_validate = split_set(sample, training_percentage=0.90)

    print("x_train", x_train.shape)
    print("y_train", y_train.shape)
    print("x_validate", x_validate.shape)
    print("y_validate", y_validate.shape)

    model = train_model(x_train, y_train, x_validate, y_validate)
    check_model(model, x_validate, y_validate)


def prepare_data(data):
    data = data.drop("event_id", axis=1)
    return data


def extract_used_ressources(data):
    event = data.values[:, 0]
    resources = data.values[:, 1]

    print(len(numpy.unique(event)))
    print(len(numpy.unique(resources)))


def split_set(data_set, training_percentage=0.75):
    event_count = data_set.shape[0]
    train_size = int(event_count * training_percentage)

    classes = data_set.pop("class")

    x_train = data_set[:train_size].values
    y_train = classes[:train_size].values

    x_validate = data_set[train_size:].values
    y_validate = classes[train_size:].values

    # We need to project our data in 3D for LSTM
    x_train = np.reshape(x_train, (x_train.shape[0], 1, x_train.shape[1]))
    x_validate = np.reshape(x_validate, (x_validate.shape[0], 1, x_validate.shape[1]))

    return x_train, y_train, x_validate, y_validate


def train_model(x_train, y_train, x_validate, y_validate):
    model = Sequential()
    model.add(LSTM(units=50, input_shape=(x_train.shape[1], x_train.shape[2])))
    # model.add(Dropout(0.2))
    # model.add(LSTM(units=50, return_sequences=True))

    model.add(Dense(units=1))
    model.compile(optimizer='adam', loss='mean_squared_error')

    model.fit(x_train, y_train, epochs=1, batch_size=72, validation_data=(x_validate, y_validate))
    history = model.fit(x_train, y_train, epochs=50, batch_size=72, validation_data=(x_validate, y_validate), verbose=2,
                        shuffle=False)
    # plot history
    pyplot.plot(history.history['loss'], label='train')
    pyplot.plot(history.history['val_loss'], label='test')
    pyplot.legend()
    pyplot.show()
    return model


def check_model(model, x_eval, y_eval):
    y_predicted = model.predict(x_eval)
    print(y_predicted.shape, y_eval.shape)
    # calculate RMSE
    diff = 0
    rmse = np.sqrt(mean_squared_error(y_predicted, y_eval))
    for i in range(len(y_predicted)):
        if y_predicted[i] != y_eval[i]:
            diff += 1

    print(diff, y_predicted.shape)
    print('Test RMSE: %.3f' % rmse)


def predict_next(model, x):
    pass


if __name__ == "__main__":
    t = 1
    predict_sequence()
