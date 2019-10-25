<template>
  <div>
    <button text="Sign in with Twitter">
      <a href="/oauth/tw/signin">
        <img
          src="https://cdn.cms-twdigitalassets.com/content/dam/developer-twitter/images/sign-in-with-twitter-gray.png"
        />
      </a>
    </button>
  </div>
</template>

<script>
import { mapGetters, mapActions } from "vuex";

export default {
  name: "SignIn",
  components: {},
  methods: {
    ...mapActions(["getUser"])
  }
};
</script>

<style></style>