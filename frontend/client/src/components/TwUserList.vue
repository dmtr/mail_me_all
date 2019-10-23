<template>
  <v-list subheader>
    <v-subheader>Users</v-subheader>
    <v-list-item v-for="user in userList" :key="user.id">
      <v-list-item-content>
        <FBUser v-bind:fbuser="user" />
      </v-list-item-content>
    </v-list-item>
  </v-list>
</template>

<script>
import TwUser from "./TwUser";

export default {
  name: "TwUserList",
  components: { TwUser },
  props: {
    userList: Array
  }
};
</script>