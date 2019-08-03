import numpy as np
import pandas as pd
from keras import Sequential
from keras.layers import Dropout, LSTM, Dense
from matplotlib import pyplot
from sklearn.decomposition import PCA
from sklearn.metrics import accuracy_score, confusion_matrix
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import MinMaxScaler


def predict_sequence():
    sample = pd.read_csv("../seq/sample.csv")
    # ressources = pd.read_csv("../seq/res.csv")
    sample = prepare_data(sample)

    # create series and split in training and validation sets
    x_series, y_series = create_series(sample, 30, 2)
    x_train, x_validate, y_train, y_validate = train_test_split(x_series, y_series, test_size=0.2)

    print("x_train", x_train.shape)
    print("y_train", y_train.shape)
    print("x_validate", x_validate.shape)
    print("y_validate", y_validate.shape)

    # Now we can train and check our model
    model = train_model(x_train, y_train, x_validate, y_validate)

    check_model(model, x_train, y_train)
    check_model(model, x_validate, y_validate)


def prepare_data(data):
    """
    Prepare the data before training
    :param data: raw data
    :return: prepared data
    """
    data = data.drop("event_id", axis=1)
    data = data.drop("timestamp", axis=1)

    prepared_data = data.values

    attr = data.values[:, 1:]
    scaler = MinMaxScaler(feature_range=(0, 1))
    attr = scaler.fit_transform(attr)

    # pca = PCA(n_components=0.90, svd_solver='full')
    # attr = pca.fit_transform(attr)

    prepared_data = np.insert(attr, 0, values=prepared_data[:, 0], axis=1)
    return prepared_data


def split_set(x, y, training_percentage=0.75):
    series_count = x.shape[0]
    split_index = int(series_count * training_percentage)

    x_train = x[:split_index]
    y_train = y[:split_index]

    x_validate = x[split_index:]
    y_validate = y[split_index:]

    return x_train, y_train, x_validate, y_validate


def create_series(data, used=1, predicted=1):
    """
    Splits the data in prediction sets as follows
    x is the vector of data used for the prediction
    y is the vector of class that happen after x

    :param data : pandas DataFrame, raw data
    :param used : number of event to use to predict
    :param predicted : number of event we want to predict
    """
    data_count, attribute_count = data.shape
    series_count = data_count - used - predicted

    x = np.zeros(shape=(series_count, used, attribute_count))
    y = np.zeros(shape=(series_count, predicted))

    for i in range(0, series_count):
        range_start = i
        used_range_end = range_start + used
        predicted_range_end = used_range_end + predicted

        x[i, :, :] = data[range_start:used_range_end]
        y[i, ] = data[used_range_end:predicted_range_end, 0]

    return x, y


def train_model(x_train, y_train, x_validate, y_validate):
    """
    Train the tf model with given sets
    :param x_train: training set attributes values
    :param y_train: training set class labels
    :param x_validate: training set attributes values
    :param y_validate: training set class labels
    :returns the trained model
    """
    model = Sequential()
    model.add(LSTM(units=64, input_shape=(x_train.shape[1], x_train.shape[2]), kernel_initializer='normal'))
    model.add(Dense(y_train.shape[1], activation='sigmoid', kernel_initializer='normal'))

    model.compile(loss='binary_crossentropy', optimizer='adam', metrics=['accuracy'])
    print(model.summary())

    history = model.fit(x_train, y_train, epochs=1, batch_size=64, validation_data=(x_validate, y_validate), verbose=2)

    # plot history
    pyplot.plot(history.history['loss'], label='train')
    pyplot.plot(history.history['val_loss'], label='test')
    pyplot.legend()
    pyplot.show()

    return model


def check_model(model, x_eval, y_eval):
    for x in x_eval:
        print(model.predict(x.reshape(1, x.shape[0], x.shape[1])))
    # y_predicted = model.predict(x_eval)

    # print(accuracy_score(y_eval, y_predicted))
    # print(confusion_matrix(y_eval, y_predicted))


if __name__ == "__main__":
    predict_sequence()
