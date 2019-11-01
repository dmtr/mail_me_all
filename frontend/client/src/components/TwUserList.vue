<template>
  <v-list>
    <v-list-item v-for="item in userList" :key="item.id">
      <v-list-item-avatar>
        <v-img :src="item.profile_image_url"></v-img>
      </v-list-item-avatar>
      <v-list-item-content>
        <a
          v-bind:href="`https://twitter.com/${item.screen_name}`"
          target="_blank"
          rel="noopener noreferrer"
        >
          <v-list-item-title v-text="item.name"></v-list-item-title>
        </a>
      </v-list-item-content>
      <v-list-item-action>
        <v-btn icon v-on:click="$emit('removeUser', item)">
          <v-icon color="grey lighten-1">mdi-account-remove</v-icon>
        </v-btn>
      </v-list-item-action>
    </v-list-item>
  </v-list>
</template>

<script>
export default {
  name: "TwUserList",
  props: {
    userList: Array
  }
};
</script>