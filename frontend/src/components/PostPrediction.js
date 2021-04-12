import React, { Component } from 'react';
import axios from 'axios';

import {
  Container, Col, Form, Row, Label,
  FormGroup, Button, FormText, FormFeedback, Input
} from 'reactstrap';

const api = axios.create({
  baseURL: window.REACT_APP_BACKEND_URL,
  headers: {
    'Content-Type': 'application/json',
  }
})

class PostPrediction extends Component {

  constructor(props) {
    super(props);

    this.state = {
      isLoading: false,
      prediction: '',
      url: '',
      validate :{
        urlState: '',
      },
    }
    this.handleChange = this.handleChange.bind(this);
    this.baseState = this.state; // preserver the initial state in a new object
  }

  resetForm = () => {
    this.setState(this.baseState)
  }

  handleChange = event => {
    const imageUrlRex = /^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/)?[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)jpg$/;
    const { validate } = this.state
    if (imageUrlRex.test(event.target.value)) {
      validate.urlState = 'has-success'
    } else {
      validate.urlState = 'has-danger'
    }

    this.setState({validate})
    this.setState({ [event.target.name]: event.target.value });
  }

  handleClick= event => {
    event.preventDefault(); //stop browser from reloading page

    const URL = {
      url: this.state.url
    }

    this.setState({isLoading: true});
    api
    .post('/classify', URL)
    .then(res => {
      this.setState({
        prediction: res.data.label,
        isLoading: false,
      })
      }
    )
    this.image = this.state.url
    this.setState({
      url: ''
    })
  }

  handleCancelClick = event => {
    this.resetForm()
  }

  capitaliseString(str) {
    return str.charAt(0).toUpperCase() + str.slice(1);
  }

  render() {
    const url = this.state.url;
    const isLoading = this.state.isLoading;
    const prediction = this.state.prediction;

    return (
      <Container>
        <div>
        <h2 className="title">Image Classifier</h2>
        </div>
        <div className="content">
          <Form>
            <Row form>
              <FormGroup>
                <Label>Provide ImageUrl</Label>
                <Input
                  name="url"
                  type="url"
                  id="exampleUrl"
                  placeholder="https://image.jpg"
                  value={url}
                  valid={ this.state.validate.urlState === 'has-success' }
                  invalid={ this.state.validate.urlState === 'has-danger' }
                  onChange={this.handleChange} />
                  <FormText className="text-muted">Image URL should be of type jpeg.</FormText>
                  <FormFeedback valid>
                    This looks like a valid URL
                  </FormFeedback>
                  <FormFeedback invalid>
                    This is an invalid URL
                  </FormFeedback>
              </FormGroup>
            </Row>

            <Row>
              <Col>
                <Button
                  block
                  color="success"
                  type="submit"
                  disabled={url === "" || this.state.validate.urlState === 'has-danger' }
                  onClick={!isLoading ? this.handleClick : null}>
                  { isLoading ? 'Making Prediction' : 'Predict' }
                </Button>
              </Col>
              <Col>
                <Button
                  block
                  color="danger"
                  disabled={isLoading}
                  onClick={this.handleCancelClick}>
                  Reset prediction
                </Button>
              </Col>
            </Row>
          </Form>
          
          {prediction === "" ? null :
            (<Row>
              <Col className="result-container">
                <h5 id="result">This looks like: {this.capitaliseString(prediction)}</h5>
                <img className="img-responsive" alt="new" src={this.image} width="300" />
              </Col>
            </Row>)
          }
        </div>
      </Container>
  );
  }
}

export default PostPrediction