<template>
  <div>
    <v-facebook-login app-id="2288197271493743" @login="login" />
  </div>
</template>

<script>
import axios from "axios";
import { mapGetters, mapActions } from "vuex";
import { VFBLogin as VFacebookLogin } from "vue-facebook-login-component";

export default {
  name: "SignIn",
  components: {
    VFacebookLogin
  },
  methods: {
    ...mapActions(["getUser"]),
    login: function(event) {
      console.debug("login");
      if (event.status !== "connected") {
        return;
      }
      const this_ = this;
      const data = event.authResponse;
      const res = axios
        .post("api/signin/fb", {
          fbid: data.userID,
          fbtoken: data.accessToken
        })
        .then(function(response) {
          this_.getUser();
        })
        .catch(function(error) {
          console.log(error);
        });
    }
  }
};
</script>

<style></style>
