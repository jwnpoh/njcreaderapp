<script>
  export let data;
  export let form;

  let note = data.data;
  let tldr = note.tldr;
  let examples = note.examples;
  let notes = note.notes;

  let tldrTags = form?.tldrTags ?? [];
  let examplesTags = form?.examplesTags ?? [];
  let notesTags = form?.notesTags ?? [];

  let tldrLength = 0;
  const updateLength = () => {
    tldrLength = tldr.length;
  };

  const parseTldrTags = (input) => {
    tldrTags = input.match(/#[a-z\d-]+/gi);
  };
  const parseExamplesTags = (input) => {
    examplesTags = input.match(/#[a-z\d-]+/gi);
  };
  const parseNotesTags = (input) => {
    notesTags = input.match(/#[a-z\d-]+/gi);
  };
</script>

<form method="POST" action="?/save">
  <input name="note_id" type="hidden" hidden class="hidden" value={note.id} />
  <input
    name="article_id"
    type="hidden"
    hidden
    class="hidden"
    value={note.article_id}
  />
  <div class="flex pt-4 md:pt-6 justify-center items-center ">
    <div
      class="card w-11/12 md:w-10/12 shadow-sm shadow-yellow-200 bg-yellow-300 bg-opacity-10"
    >
      <div class="card-body px-6 md:px-12">
        <div class="lg:inline-flex">
          <h1 class="card-title min-w-fit flex-1 pb-2 text-3xl underline">
            Edit note
          </h1>
          <div
            class="alert alert-success bg-opacity-50 flex-shrink lg:mx-10 py-3"
          >
            <p class="text-sm">
              Tip: add #hashtags to organise your notes for future reference.
              You can add hashtags anywhere in the form, even in the #middle of
              your note!
            </p>
          </div>
        </div>

        <input
          name="article_title"
          type="hidden"
          hidden
          class="hidden"
          value={note.article_title}
        />
        <input
          name="article_url"
          type="hidden"
          hidden
          class="hidden"
          value={note.article_url}
        />
        <div class="pt-2">
          <h2 class="text-lg font-semibold">
            <a href={note.article_url} target="_blank" rel="noreferrer"
              >{note.article_title}</a
            >
          </h2>
        </div>

        <input
          name="tldr_tags"
          bind:value={tldrTags}
          type="hidden"
          hidden
          class="hidden"
        />
        <div class="pt-4">
          <h2 class="text-lg font-semibold">tl;dr</h2>
          <p class="pb-2 text-md  text-justify opacity-80">
            Write a short summary of your main takeaways from reading this
            article. <br />You are encouraged to make use of critical thinking
            tools such as:
          </p>
          <ul class="pl-5 pb-2 text-md opacity-80">
            <li class="list-disc list-outside">
              <strong>Mental Models</strong> (problem, stakeholders, causes, consequences,
              solutions, implications)
            </li>
            <li class="list-disc list-outside">
              <strong>Paul's Elements of Reasoning</strong> (issue, purpose, point
              of view, assumptions, concepts, evidence, inferences, implications).
            </li>
          </ul>

          <textarea
            required
            maxlength="160"
            name="tldr"
            type="text"
            placeholder="My main takeaways..."
            class="textarea px-2 py-2 w-screen max-w-full bg-secondary bg-opacity-5"
            bind:value={tldr}
            on:input={updateLength}
            on:input={parseTldrTags(tldr)}
          />
          <label for="tldr" class="label">
            <span class="label-text-alt" />
            <span class="label-text-alt">{tldrLength}/160</span>
          </label>
        </div>

        <input
          name="examples_tags"
          bind:value={examplesTags}
          type="hidden"
          hidden
          class="hidden"
        />
        <div>
          <h2 class="text-lg font-semibold">Examples</h2>
          <p class="pb-2 text-md  text-justify opacity-80">
            Note down any useful or interesting examples that you find in the
            article.
          </p>
          <textarea
            required
            maxlength="160"
            name="examples"
            type="text"
            placeholder="Interesting or useful examples..."
            class="textarea px-2 py-2 w-screen max-w-full bg-secondary bg-opacity-5"
            bind:value={examples}
            on:input={parseExamplesTags(examples)}
          />
        </div>

        <input
          name="notes_tags"
          bind:value={notesTags}
          type="hidden"
          hidden
          class="hidden"
        />
        <div class="pt-4">
          <h2 class="text-lg font-semibold">Further notes</h2>
          <p class="pb-2 text-md  text-justify opacity-80">
            Note down any further reflection that you have. You can consider the
            following prompts:
          </p>
          <ul class="pl-5 text-md opacity-80">
            <li class="list-disc list-outside">
              How does this article relate to myself/my society?
            </li>
            <li class="list-disc list-outside">
              How might I use the content from this article to answer a
              particular past year question?
            </li>
            <li class="list-disc list-outside">
              Does this article remind me of something else I have read or
              watched?
            </li>
          </ul>
          <textarea
            name="notes"
            type="text"
            placeholder="Further reflections..."
            class="textarea my-2 px-2 py-2 w-screen max-w-full bg-secondary bg-opacity-5"
            bind:value={notes}
            on:input={parseNotesTags(notes)}
          />
        </div>

        <div class="flex py-2 ">
          <label class="px-2" for="make_public">Make note public?</label>
          <input
            name="make_public"
            type="checkbox"
            class="checkbox"
            bind:checked={note.public}
          />
        </div>
        <div class="py-2">
          <button class="btn btn-sm btn-secondary "
            >Save changes to notebook</button
          >
        </div>
        {#if form?.error}
          <p class="text-primary">{form?.message}</p>
        {/if}
        {#if form?.success}
          <div
            class="alert alert-success z-50 fixed place-self-center top-1/3 max-w-sm shadow-2xl border-2 border-black border-opacity-40"
          >
            <div class="grid md:grid-cols-2 p-5 place-items-center">
              <p class="text-xl col-span-2 text-center">
                Successfully added post! <br /> What would you like to do now?
              </p>
              <span class="col-span-2" />
              <a href="/">Go to article feed</a>
              <a href="/notes">Go to notebook</a>
            </div>
          </div>
        {/if}
      </div>
    </div>
  </div>
</form>
