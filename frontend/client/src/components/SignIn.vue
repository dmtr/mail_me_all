<template>
  <div>
    <v-facebook-login app-id="2288197271493743" v-on:login="login" v-on:logout="logout"></v-facebook-login>
  </div>
</template>

<script>
import axios from "axios";
import { VFBLogin as VFacebookLogin } from "vue-facebook-login-component";

export default {
  name: "SignIn",
  components: {
    VFacebookLogin
  },
  methods: {
    login: function(event) {
      console.debug("login");
      if (event.status !== "connected") {
        return;
      }

      const data = event.authResponse;
      const res = axios
        .post("api/signin/fb", {
          fbid: data.userID,
          fbtoken: data.accessToken
        })
        .then(function(response) {
          // handle success
          console.log(response);
        })
        .catch(function(error) {
          // handle error
          console.log(error);
        })
        .finally(function() {
          // always executed
        });
      console.debug(res);
    },
    logout: function(event) {
      console.debug("logout");
    }
  }
};
</script>

<style></style>
