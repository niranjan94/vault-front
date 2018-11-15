<template>
  <div>
    <h4 class="ui top attached header">
      Word Cloud
    </h4>
    <div class="ui attached segment" :class="{ loading: isLoading }">
      <form class="ui small form" v-on:submit.prevent="addWord">
        <div class="fields">
          <div class="seven wide required field">
            <label>Word</label>
            <input type="text" v-model.trim="target" placeholder="" required>
          </div>
          <div class="seven wide field">
            <label>Replacement <small class="muted">(optional)</small></label>
            <input type="text" v-model.trim="replacement" placeholder="">
          </div>
          <div class="two wide field">
            <label>&nbsp;</label>
            <button type="submit" class="ui fluid button" :disabled="isSaving">Add</button>
          </div>
        </div>
      </form>
    </div>
    <div class="ui bottom attached segment" v-if="words.length > 0">
      <table class="ui very compact very basic small table">
        <thead>
        <tr>
          <th class="one wide">#</th>
          <th class="six wide">Word</th>
          <th class="six wide">Replacement</th>
          <th class="one wide"> </th>
        </tr>
        </thead>
        <tbody>

        <tr v-for="(word, index) in words">
          <td>{{index + 1}}</td>
          <td>{{word.target}}</td>
          <td>{{!word.replacement || word.replacement.trim() === '' ? '-' : word.replacement}}</td>
          <td>
            <button class="ui mini red icon button"
                    v-on:click.prevent="deleteWord(word)"
                    :disabled="isSaving"
                    data-tooltip="Delete"
                    data-position="right center"
                    data-inverted="">
              <i class="trash icon"></i>
            </button>
          </td>
        </tr>
        </tbody>
      </table>

    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import axios from 'axios';

@Component({})
export default class WordCloud extends Vue {

  private target: string = '';
  private replacement: string = '';

  private isSaving: boolean = false;
  private isLoading: boolean = true;

  private words: any[] = [];

  protected async created() {
    this.isLoading = true;
    try {
      const response = await axios.get('words');
      this.words = response.data;
    } catch (e) {
      this.words = [];
      this.$notify({ type: 'error', text: e.message });
    }
    this.isLoading = false;
  }

  protected async addWord() {
    await this.save(
      this.words.concat([{ target: this.target, replacement: this.replacement }])
    );
    this.target = this.replacement = '';
    this.$notify({ type: 'success',  text: 'The word has been added.' });
  }

  protected async deleteWord(word) {
    await this.save(
      this.words.filter((it) => it !== word)
    );
    this.$notify({ type: 'success',  text: 'The word has been deleted.' });
  }

  private async save(words) {
    this.isSaving = true;
    try {
      await axios.post('words', words);
      this.words = words;
    } catch (e) {
      this.$notify({ type: 'error', text: e.message });
    }
    this.isSaving = false;
  }
}
</script>

<style scoped lang="scss">
  .muted {
    color: gray;
  }
</style>
