import axios from "axios";
import _ from "lodash";

const defaultErorr = "Server Error";

const getError = e => ({
  code: e.code,
  message: _.get(e, "message", defaultErorr)
});

class ApiResult {
  constructor(data, error) {
    this.data = data;
    this.error = error ? getError(error) : error;
  }
}

export async function getUser() {
  const errors = [401, 404, 500];
  try {
    const response = await axios.get("api/user");
    return new ApiResult(response.data, null);
  } catch (error) {
    if (error.response && errors.indexOf(error.response.status) != -1) {
      return new ApiResult({ signedIn: false, name: "", id: "" }, null);
    } else {
      console.log(error);
      return new ApiResult(null, error);
    }
  }
}

export async function getSubscriptions(userId) {
  try {
    const response = await axios.get(`api/user/${userId}/subscriptions`);
    return new ApiResult(response.data["subscriptions"], null);
  } catch (error) {
    console.log(error);
    return new ApiResult(null, error);
  }
}

export async function createSubscription(userId, subscription) {
  try {
    console.log(userId, subscription);
    const response = await axios.post(
      `api/user/${userId}/subscriptions`,
      subscription
    );
    return new ApiResult(response.data, null);
  } catch (error) {
    console.log(error);
    return new ApiResult(null, error);
  }
}
