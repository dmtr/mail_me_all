const state = {
  user: { signedIn: false },
  userLists: []
};

const getters = {
  isUserSignedIn: state => state.user.signedIn
};

export default {
  state,
  getters
};
