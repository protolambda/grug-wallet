
function onError(error) {
  console.log(error);
}

let inputendpointName = document.querySelector('#new-endpoint #name');
let inputEndpointValue = document.querySelector('#new-endpoint #endpoint');

let endpointList = document.querySelector('#endpoint-list');
let providerList = document.querySelector('#provider-list');


let addBtn = document.querySelector('.add');

addBtn.addEventListener('click', addEndpoint);

initialize();

function initialize() {
  let gettingAllStorageItems = browser.storage.local.get(null);
  gettingAllStorageItems.then((results) => {
    let serverEndpoints = Object.keys(results);
    for (let endpointKey of serverEndpoints) {
      if(endpointKey.startsWith('endpoint-')) {
        displayEndpoint(endpointKey, results[endpointKey]);
      }
    }
  }, onError);
  updateProvidersList();
}

// Add a endpointEntry to the display, and storage

function storeEndpoint(name, endpoint) {
  let storingEndpointEntry = browser.storage.local.set({ [name] : endpoint });
  storingEndpointEntry.then(() => {
    displayEndpoint(name, endpoint);
  }, onError);
}

function addEndpoint() {
  let endpointKey = 'endpoint-' + inputendpointName.value;
  let endpointURL = inputEndpointValue.value.trim();
  let gettingItem = browser.storage.local.get(endpointKey);
  gettingItem.then((result) => {
    let objTest = Object.keys(result);
    if(objTest.length < 1 && endpointKey !== '' && endpointURL !== '') {
      inputendpointName.value = '';
      inputEndpointValue.value = '';
      storeEndpoint(endpointKey, endpointURL);
    }
  }, onError);
}

// function to show endpoint
function displayEndpoint(endpointKey, endpointURL) {
  if (!endpointKey.startsWith("endpoint-")) { // endpoints must have prefixed names internally
    return
  }

  // create endpointEntry display box
  let endpointEntry = document.createElement('div');
  let endpointEntryDisplay = document.createElement('div');
  let endpointEntryH = document.createElement('h2');
  let endpointEntryPara = document.createElement('p');
  let deleteBtn = document.createElement('button');
  let activateBtn = document.createElement('button');
  let clearFix = document.createElement('div');

  endpointEntry.setAttribute('class','endpointEntry');

  endpointDisplayName = endpointKey.substring("endpoint-".length);
  endpointEntryH.textContent = endpointDisplayName;
  endpointEntryPara.textContent = endpointURL;
  deleteBtn.setAttribute('class','delete');
  deleteBtn.textContent = 'Delete endpoint';
  activateBtn.setAttribute('class','inject');
  activateBtn.textContent = 'Load endpoint';
  clearFix.setAttribute('class','clearfix');

  endpointEntryDisplay.appendChild(endpointEntryH);
  endpointEntryDisplay.appendChild(endpointEntryPara);
  endpointEntryDisplay.appendChild(deleteBtn);
  endpointEntryDisplay.appendChild(activateBtn);
  endpointEntryDisplay.appendChild(clearFix);

  endpointEntry.appendChild(endpointEntryDisplay);

  // delete option
  deleteBtn.addEventListener('click',(e) => {
    const evtTarget = e.target;
    evtTarget.parentNode.parentNode.parentNode.removeChild(evtTarget.parentNode.parentNode);
    browser.storage.local.remove(endpointKey);
  })

  activateBtn.addEventListener('click',async (e) => {
    await browser.runtime.sendMessage({command: 'loadProviders', server: endpointURL});
  })

  // edit box
  let endpointEntryEdit = document.createElement('div');
  let endpointEntryTitleEdit = document.createElement('input');
  let endpointEntryBodyEdit = document.createElement('input');
  let clearFix2 = document.createElement('div');

  let updateBtn = document.createElement('button');
  let cancelBtn = document.createElement('button');

  updateBtn.setAttribute('class','update');
  updateBtn.textContent = 'Update endpointEntry';
  cancelBtn.setAttribute('class','cancel');
  cancelBtn.textContent = 'Cancel update';

  endpointEntryEdit.appendChild(endpointEntryTitleEdit);
  endpointEntryTitleEdit.value = endpointDisplayName;
  endpointEntryEdit.appendChild(endpointEntryBodyEdit);
  endpointEntryBodyEdit.value = endpointURL;
  endpointEntryEdit.appendChild(updateBtn);
  endpointEntryEdit.appendChild(cancelBtn);

  endpointEntryEdit.appendChild(clearFix2);
  clearFix2.setAttribute('class','clearfix');

  endpointEntry.appendChild(endpointEntryEdit);

  endpointList.appendChild(endpointEntry);
  endpointEntryEdit.style.display = 'none';

  endpointEntryH.addEventListener('click',() => {
    endpointEntryDisplay.style.display = 'none';
    endpointEntryEdit.style.display = 'block';
  })

  endpointEntryPara.addEventListener('click',() => {
    endpointEntryDisplay.style.display = 'none';
    endpointEntryEdit.style.display = 'block';
  })

  cancelBtn.addEventListener('click',() => {
    endpointEntryDisplay.style.display = 'block';
    endpointEntryEdit.style.display = 'none';
    endpointEntryTitleEdit.value = endpointDisplayName;
    endpointEntryBodyEdit.value = endpointURL;
  })

  updateBtn.addEventListener('click',() => {
    if(endpointEntryTitleEdit.value !== endpointKey || endpointEntryBodyEdit.value !== endpointURL) {
      updateEndpoint(endpointKey,"endpoint-"+endpointEntryTitleEdit.value, endpointEntryBodyEdit.value);
      endpointEntry.parentNode.removeChild(endpointEntry);
    }
  });
}

function updateEndpoint(delEndpoint,newendpointKey,newEndpointURL) {
  let storingEndpointEntry = browser.storage.local.set({ [newendpointKey] : newEndpointURL });
  storingEndpointEntry.then(() => {
    if(delEndpoint !== newendpointKey) {
      let removingEndpointEntry = browser.storage.local.remove(delEndpoint);
      removingEndpointEntry.then(() => {
        displayEndpoint(newendpointKey, newEndpointURL);
      }, onError);
    } else {
      displayEndpoint(newendpointKey, newEndpointURL);
    }
  }, onError);
}

async function updateProvidersList() {

  const providerKeys = Object.keys(await browser.storage.local
      .get(null)).filter((key) => key.startsWith("web3-"));

  const getProviderElem = async (providerKey) => {
    let providerElem = document.createElement('div');

    const providerEndpointDef = (await browser.storage.local.get(providerKey))[providerKey];
    if (!providerEndpointDef) {
      console.log("missing provider entry for ", providerKey);
      return providerElem;
    }

    let providerContent = document.createElement('div');
    const info = providerEndpointDef.info;
    providerContent.innerText = "name: "+info.name+" uuid: "+info.uuid+" endpoint: "+providerEndpointDef.endpoint;

    let providerIcon = document.createElement('img');
    providerIcon.src = info.icon;
    providerContent.appendChild(providerIcon);

    providerElem.appendChild(providerContent);

    let providerErr = document.createElement('div');
    providerErr.appendChild(providerElem);

    let announceBtn = document.createElement('button');
    announceBtn.setAttribute('class','announce');
    announceBtn.textContent = 'announce';
    announceBtn.addEventListener('click',async () => {
      await browser.tabs.query({active: true, currentWindow: true})
          .then(async (tabs) => {
            await browser.tabs.sendMessage(tabs[0].id, {
              command: "shareProvider",
              providerUUID: info.uuid,
            });
          })
          .catch((err) => {
            providerErr.innerText = 'fail: '+err.message;
          });
    })
    providerElem.appendChild(announceBtn);

    let overrideGlobalBtn = document.createElement('button');
    overrideGlobalBtn.setAttribute('class','overrideEth');
    overrideGlobalBtn.textContent = 'override';
    overrideGlobalBtn.addEventListener('click',async () => {
      await browser.tabs.query({active: true, currentWindow: true})
          .then(async (tabs) => {
            await browser.tabs.sendMessage(tabs[0].id, {
              command: "overrideGlobalProvider",
              providerUUID: info.uuid,
            });
          })
          .catch((err) => {
            providerErr.innerText = 'fail: '+err.message;
          });
    })
    providerElem.appendChild(overrideGlobalBtn);
    return providerElem;
  }

  providerList.innerHTML = "";
  await Promise.all(providerKeys.map(async (provKey) => {
    const providerElem = await getProviderElem(provKey);
    providerList.appendChild(providerElem);
  }));

}

browser.runtime.onMessage.addListener(async (message) => {
  if (message.command === "refreshProviders") {
    console.log("refreshing providers");
    await updateProvidersList();
  }
})

