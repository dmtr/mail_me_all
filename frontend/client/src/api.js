import axios from "axios";
import _ from "lodash";

const defaultErorr = "Server Error";

function getError(e) {
  if (e.response) {
    const data = e.response.data;
    if (data && typeof data === "object") {
      return {
        code: data.code,
        message: data.message
      };
    } else {
      return {
        code: e.status,
        message: e.message
      };
    }
  } else {
    return {
      code: null,
      message: _.get(e, "message", defaultErorr)
    };
  }
}

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
      console.log(error.toJSON());
      return new ApiResult(null, error);
    }
  }
}

export async function getSubscriptions() {
  try {
    const response = await axios.get("api/subscriptions");
    return new ApiResult(response.data["subscriptions"], null);
  } catch (error) {
    console.log(error.toJSON());
    return new ApiResult(null, error);
  }
}

export async function createSubscription(subscription) {
  try {
    const response = await axios.post("api/subscriptions", subscription);
    return new ApiResult(response.data, null);
  } catch (error) {
    console.log(error.toJSON());
    return new ApiResult(null, error);
  }
}

export async function updateSubscription(subscription) {
  try {
    const response = await axios.put("api/subscriptions", subscription);
    return new ApiResult(response.data, null);
  } catch (error) {
    console.log(error.toJSON());
    return new ApiResult(null, error);
  }
}

export async function deleteSubscription(subscriptionId) {
  try {
    const response = await axios.delete(`api/subscriptions/${subscriptionId}`);
    return new ApiResult(response.data, null);
  } catch (error) {
    console.log(error.toJSON());
    return new ApiResult(null, error);
  }
}

export async function deleteAccount() {
  try {
    const response = await axios.delete(`api/user`);
    return new ApiResult(response.data, null);
  } catch (error) {
    console.log(error.toJSON());
    return new ApiResult(null, error);
  }
}
