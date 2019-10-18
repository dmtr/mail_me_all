import axios from "axios";

const state = {
  user: { signedIn: false },
  userLists: []
};

const getters = {
  isUserSignedIn: state => state.user.signedIn
};

const actions = {
  getUser({ commit, state }) {
    const res = axios
      .get("api/user")
      .then(function(response) {
        console.log(response.data);
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
