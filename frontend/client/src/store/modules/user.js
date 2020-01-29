import _ from "lodash";

import {
  getUser,
  getSubscriptions,
  createSubscription,
  updateSubscription,
  deleteSubscription,
  deleteAccount
} from "../../api";

const state = {
  user: null,
  subscriptions: []
};

const getters = {
  isUserLoaded: state => (state.user ? true : false),
  isUserSignedIn: state => state.user && state.user.signedIn,
  subscriptions: state => state.subscriptions,
  email: state =>
    state.subscriptions.length
      ? _.chain(state.subscriptions)
          .groupBy("email")
          .sortBy("length")
          .last()
          .value()[0]["email"]
      : ""
};

function handle401(commit, res) {
  if ((res.error && res.error.code === 5) || res.error.code == 401) {
    commit("setUser", null);
  }
}

const actions = {
  async getUser({ commit }) {
    const res = await getUser();
    if (!res.error) {
      commit("setUser", res.data);
    }
    return res;
  },

  async createSubscription(context, subscription) {
    const res = await createSubscription(subscription);
    handle401(context.commit, res);
    return res;
  },

  async updateSubscription(context, subscription) {
    const res = await updateSubscription(subscription);
    handle401(context.commit, res);
    return res;
  },

  async deleteSubscription({ commit }, subscription) {
    const res = await deleteSubscription(subscription.id);
    if (!res.error) {
      commit("removeSubscription", subscription);
    } else {
      handle401(commit, res);
    }
  },

  async getSubscriptions({ commit }) {
    const res = await getSubscriptions();
    if (!res.error) {
      commit("setSubscriptions", res.data);
    } else {
      handle401(commit, res);
    }
    return res;
  },

  async deleteAccount({ commit }) {
    const res = await deleteAccount();
    if (!res.error) {
      commit("setUser", null);
      commit("setSubscriptions", []);
    } else {
      handle401(commit, res);
    }
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
  },

  removeSubscription(state, subscription) {
    state.subscriptions = _.filter(
      state.subscriptions,
      s => s.id != subscription.id
    );
  }
};

export default {
  state,
  getters,
  mutations,
  actions
};
