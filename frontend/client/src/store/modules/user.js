import axios from "axios";

const state = {
  user: null,
  subscriptionsList: [
    {
      id: 1,
      title: "foo",
      day: "monday",
      email: "foo@bar.com",
      userList: [{ id: "1", name: "john" }]
    }
  ]
};

const getters = {
  isUserLoaded: state => (state.user ? true : false),
  isUserSignedIn: state => state.user && state.user.signedIn,
  subscriptionsList: state => state.subscriptionsList
};

const actions = {
  getUser({ commit, state }) {
    const res = axios
      .get("api/user")
      .then(function(response) {
        commit("setUser", response.data);
      })
      .catch(function(error) {
        console.log(error);
      });
  }
};

const mutations = {
  setUser(state, user) {
    state.user = user;
  }
};

export default {
  state,
  getters,
  mutations,
  actions
};
