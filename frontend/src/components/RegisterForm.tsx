import { Button, Col, Container, Row, Stack } from "react-bootstrap";
import Form from "react-bootstrap/Form";
import InputGroup from 'react-bootstrap/InputGroup';
import { Link, NavigateFunction, useNavigate } from "react-router-dom";
import AccountService from "../service/AccountService";
import { Formik } from "formik";

import * as yup from "yup";


const RegisterForm: React.FC = () => {
    let navigate: NavigateFunction = useNavigate();

    const schema = yup.object().shape({
        name: yup.string().required("name can't be empty."),
        username: yup.string().required("username can't be empty."),
        email: yup.string().required("email can't be empty."),
        password: yup.string().required("password  can't be empty."),
        repeatPassword: yup.string().required("passwords do not match."),
    });


    const handleSubmit = (values: any, helpers: any) => {

        console.log(values)

        AccountService.register(values.name, values.username, values.email, values.password).then(
            () => {
                navigate("/");
                window.location.reload()
            },
            (error) => {
                const resMessage =
                    (error.response &&
                        error.response.data &&
                        error.response.data.message) ||
                    error.message ||
                    error.toString();

                if (error.code === "ERR_BAD_REQUEST") {
                    helpers.setFieldError("username", error)
                }
            }
        );

    }

    return (
        <Container className="my-4 w-50">
            <Formik
                validationSchema={schema}
                validateOnChange={false}
                onSubmit={(values, helpers) => {
                    handleSubmit(values, helpers)
                }}
                initialValues={{
                    name: '',
                    username: '',
                    email: '',
                    password: '',
                    repeatPassword: '',
                }}
            >
                {({
                    handleSubmit,
                    handleChange,
                    values,
                    touched,
                    errors,
                }) => (
                    <Form noValidate onSubmit={handleSubmit}>
                        <Row className="mb-4">
                            <Row>
                                <Col className="h-100 d-flex justify-content-center align-items-center">
                                    <h2>Sign in</h2>
                                </Col>
                            </Row>
                        </Row>
                        <Row className="my-2">
                            <Row>
                                <Form.Group as={Col} mb="2">
                                    <Form.Control
                                        type="text"
                                        placeholder="name"
                                        name="name"
                                        value={values.name}
                                        onChange={handleChange}
                                        isInvalid={!!errors.name}
                                    />
                                    <Form.Control.Feedback type="invalid">
                                        {errors.name}
                                    </Form.Control.Feedback>
                                </Form.Group>

                                <Form.Group as={Col} mb="2" controlId="validationFormikUsername">
                                    <InputGroup hasValidation>
                                        <Form.Control
                                            type="text"
                                            placeholder="username"
                                            aria-describedby="inputGroupPrepend"
                                            name="username"
                                            value={values.username}
                                            onChange={handleChange}
                                            isInvalid={!!errors.username}
                                        />
                                        <Form.Control.Feedback type="invalid">
                                            {errors.username}
                                        </Form.Control.Feedback>
                                    </InputGroup>
                                </Form.Group>
                            </Row>
                            <Row className="my-2">
                                <Form.Group as={Col} mb="2">
                                    <Form.Control
                                        type="text"
                                        placeholder="email"
                                        name="email"
                                        value={values.email}
                                        onChange={handleChange}
                                        isInvalid={!!errors.email}
                                    />
                                    <Form.Control.Feedback type="invalid">
                                        {errors.email}
                                    </Form.Control.Feedback>
                                </Form.Group>
                            </Row>
                            <Row className="my-2">
                                <Form.Group as={Col} my="4">
                                    <Form.Control
                                        type="password"
                                        placeholder="password"
                                        name="password"
                                        value={values.password}
                                        onChange={handleChange}
                                        isInvalid={!!errors.password}
                                    />
                                    <Form.Control.Feedback type="invalid">
                                        {errors.password}
                                    </Form.Control.Feedback>
                                </Form.Group>

                                <Form.Group as={Col} my="4">
                                    <Form.Control
                                        type="password"
                                        placeholder="repeatPassword"
                                        name="repeatPassword"
                                        value={values.repeatPassword}
                                        onChange={handleChange}
                                        isInvalid={!!errors.repeatPassword}
                                    />
                                    <Form.Control.Feedback type="invalid">
                                        {errors.repeatPassword}
                                    </Form.Control.Feedback>
                                </Form.Group>
                            </Row>
                        </Row>
                        <Row>
                            <Col className="h-100 d-flex justify-content-center align-items-center">
                                <Stack gap={4} direction="horizontal">
                                    <Button type="submit">sign up</Button>
                                    <Link to={"/login"}>
                                        <Button type="submit">I already have an account</Button>
                                    </Link>
                                </Stack>
                            </Col>
                        </Row>
                    </Form>
                )}
            </Formik>
        </Container>
    )
}

export default RegisterForm;