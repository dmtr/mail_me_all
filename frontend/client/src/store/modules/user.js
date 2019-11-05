import axios from "axios";
import _ from "lodash";

const state = {
  user: null,
  subscriptions: []
};

const getters = {
  isUserLoaded: state => (state.user ? true : false),
  isUserSignedIn: state => state.user && state.user.signedIn,
  subscriptions: state => state.subscriptions
};

const actions = {
  getUser({ commit, state }) {
    const errors = [401, 404, 500];

    const res = axios
      .get("api/user")
      .then(function(response) {
        commit("setUser", response.data);
      })
      .catch(function(error) {
        if (error.response && errors.indexOf(error.response.status) != -1) {
          commit("setUser", { signedIn: false, name: "", id: "" });
        } else {
          console.log(error.message);
        }
      });
  },

  createSubscription({ commit, state }, subscription) {
    const res = axios
      .post(`api/user/${state.user.id}/subscriptions`)
      .then(function(response) {
        commit("addSubscription", response.data);
      })
      .catch(function(error) {
        console.log(error);
      });
  }
};

const mutations = {
  setUser(state, user) {
    state.user = user;
  },

  addSubscription(state, subscription) {
    if (
      subscription &&
      -1 === _.indexOf(state.subscriptions, ["id", subscription.id])
    ) {
      state.subscriptions.push(subscription);
    }
  }
};

export default {
  state,
  getters,
  mutations,
  actions
};
