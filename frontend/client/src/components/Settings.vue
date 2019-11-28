<template>
  <div v-if="isUserSignedIn">
    <v-list max-width="400px" dense>
      <v-list-item>
        <v-list-item-content>
          <v-list-item-title>Delete Account</v-list-item-title>
        </v-list-item-content>
        <v-list-item-action>
          <v-btn outlined fab color="red" @click="deleteDialog=true">
            <v-icon>mdi-account-remove</v-icon>
          </v-btn>
        </v-list-item-action>
      </v-list-item>
    </v-list>
    <v-dialog v-model="deleteDialog" max-width="500px">
      <v-card class="pa-md-4 mx-md-auto">
        <v-card-actions>
          Delete Account?
          <v-spacer></v-spacer>
          <v-btn text color="primary" @click="remove()">Delete</v-btn>
          <v-btn text color="primary" @click="deleteDialog=false">Cancel</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import { mapGetters, mapActions } from "vuex";

export default {
  name: "Settings",
  data: () => ({
    deleteDialog: false
  }),
  computed: {
    ...mapGetters(["isUserSignedIn"])
  },
  methods: {
    ...mapActions(["deleteAccount", "getUser"]),
    remove: async function() {
      await this.deleteAccount();
      this.deleteDialog = false;
      await this.getUser();
      this.$router.push("/");
    }
  }
};
</script>