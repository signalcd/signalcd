import {html, LitElement} from 'lit-element';
import {ApiClient, DeploymentApi, PipelineApi} from '../../api/client/javascript/src/index.js';

class App extends LitElement {
    static get properties() {
        return {
            deployments: {type: Array},
            pipelines: {type: Array},
        };
    }

    constructor() {
        super();

        this.deployments = [];
        this.pipelines = [];
    }

    render() {
        return html`
            <div>
                <h1>SignalCD is alive again!</h1>
                <h3>Deployments</h3>
                <ul>
                    ${this.deployments.map((d) => html`
                        <li>${d.number} - ${d.created}</li>
                    `)}
                </ul>
                <h3>Pipelines</h3>
                <ul>
                    ${this.pipelines.map((p) => html`
                        <li>${p.id} - ${p.created}</li>
                    `)}
                </ul>
            </div>
        `;
    }

    firstUpdated(changedProperties) {
        let client = new ApiClient();
        client.basePath = `${window.location.protocol}//${window.location.host}/api/v1`;

        new DeploymentApi(client).listDeployments().then((deployments) => this.deployments = deployments);
        new PipelineApi(client).listPipelines().then((pipelines) => this.pipelines = pipelines);
    }
}

customElements.define("signalcd-app", App);
