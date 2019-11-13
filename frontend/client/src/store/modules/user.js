import _ from "lodash";
import {
  getUser,
  getSubscriptions,
  createSubscription,
  updateSubscription
} from "../../api";

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
  async getUser({ commit }) {
    const res = await getUser();
    if (!res.error) {
      commit("setUser", res.data);
    }
    return res;
  },

  async createSubscription({ commit }, subscription) {
    const res = await createSubscription(subscription);
    return res;
  },

  async updateSubscription({ commit }, subscription) {
    const res = await updateSubscription(subscription);
    return res;
  },

  async getSubscriptions({ commit }) {
    const res = await getSubscriptions();
    if (!res.error) {
      commit("setSubscriptions", res.data);
    }
    return res;
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
  },

  setSubscriptions(state, subscriptions) {
    state.subscriptions = subscriptions;
  }
};

export default {
  state,
  getters,
  mutations,
  actions
};
