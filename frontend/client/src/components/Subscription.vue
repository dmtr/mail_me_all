<template>
  <v-card>
    <v-container>
      <v-row>
        {{ subscription.title }}
        <br />
        {{ subscription.day }}
        <br />
        {{ subscription.email }}
        <br />
      </v-row>
      <v-row>
        <FBUserList v-bind:user-list="subscription.userList"></FBUserList>
      </v-row>
    </v-container>
  </v-card>
</template>

<script>
import FBUserList from "./FBUserList";

export default {
  name: "Subscription",
  components: {
    FBUserList
  },
  props: {
    subscription: Object
  }
};
</script>

<style></style>
