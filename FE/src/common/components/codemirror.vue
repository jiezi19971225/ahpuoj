<template lang="pug">
    codemirror(v-model="mycode", :options="cmOptions", @input="onCmCodeChange")
</template>

<script>
// language
import 'codemirror/mode/clike/clike';
import 'codemirror/mode/go/go';
import 'codemirror/mode/python/python';
// theme css
import 'codemirror/theme/monokai.css';
// auto close brackets
import 'codemirror/addon/edit/closebrackets';
// require active-line
import 'codemirror/addon/selection/active-line';
// styleSelectedText
import 'codemirror/addon/selection/mark-selection';
import 'codemirror/addon/search/searchcursor';
// highlightSelectionMatches
import 'codemirror/addon/scroll/annotatescrollbar';
import 'codemirror/addon/search/matchesonscrollbar';

import 'codemirror/addon/search/match-highlighter';
// keyMap
import 'codemirror/addon/edit/matchbrackets';
import 'codemirror/addon/comment/comment';
import 'codemirror/addon/dialog/dialog';
import 'codemirror/addon/dialog/dialog.css';

import 'codemirror/addon/search/search';
import 'codemirror/keymap/sublime';
// foldGutter
import 'codemirror/addon/fold/foldgutter.css';
import 'codemirror/addon/fold/brace-fold';
import 'codemirror/addon/fold/comment-fold';
import 'codemirror/addon/fold/foldcode';
import 'codemirror/addon/fold/foldgutter';
import 'codemirror/addon/fold/indent-fold';
import 'codemirror/addon/fold/markdown-fold';
import 'codemirror/addon/fold/xml-fold';

export default {
  props: {
    code: {
      type: String,
      default: '',
    },
    language: {
      type: Number,
      default: 1,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      map: new Map([
        [0, 'text/x-csrc'],
        [1, 'text/x-c++src'],
        [3, 'text/x-java'],
        [6, 'text/x-python'],
        [17, 'text/x-go'],
      ]),
      cmOptions: {
        tabSize: 4,
        styleActiveLine: false,
        lineNumbers: true,
        styleSelectedText: false,
        autoCloseBrackets: true,
        line: true,
        foldGutter: true,
        readOnly: false,
        gutters: ['CodeMirror-linenumbers', 'CodeMirror-foldgutter'],
        highlightSelectionMatches: { showToken: /\w/, annotateScrollbar: true },
        mode: 'text/x-csrc',
        keyMap: 'sublime',
        matchBrackets: true,
        showCursorWhenSelecting: true,
        theme: 'monokai',
        extraKeys: { Ctrl: 'autocomplete' },
      },
      mycode: this.code,
    };
  },
  watch: {
    language() {
      console.log(this.map.get(this.language));
      this.cmOptions.mode = this.map.get(this.language);
      console.log(this.cmOptions.mode);
    },
  },
  mounted() {
    this.cmOptions.readOnly = this.readonly;
    setTimeout(() => {
      this.styleSelectedText = true;
      this.styleActiveLine = true;
    }, 1800);
  },
  methods: {
    onCmCodeChange(newCode, event) {
      this.$emit('update:code', newCode);
    },
  },
};
</script>
<style lang="scss">
.cm-matchhighlight {
  background-color: rgba(lightblue, 0.4);
}

.CodeMirror-focused .cm-matchhighlight {
  background-image: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAIAAAACCAYAAABytg0kAAAAFklEQVQI12NgYGBgkKzc8x9CMDAwAAAmhwSbidEoSQAAAABJRU5ErkJggg==);
  background-position: bottom;
  background-repeat: repeat-x;
}

.CodeMirror-selection-highlight-scrollbar {
  background-color: green;
}

.CodeMirror {
  line-height: 20px;
  height: 550px;
}

.CodeMirror-sizer {
  margin-left: 50px !important;
}
</style>
