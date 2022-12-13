import { Button, Col, Container, Row, Stack } from "react-bootstrap";
import Form from "react-bootstrap/Form";
import InputGroup from 'react-bootstrap/InputGroup';
import { Link, NavigateFunction, useNavigate } from "react-router-dom";
import AccountService from "../service/AccountService";
import { Formik } from "formik";

import * as yup from "yup";

const LoginForm: React.FC = () => {

    let navigate: NavigateFunction = useNavigate();

    const schema = yup.object().shape({
        username: yup.string().required("username can't be empty."),
        password: yup.string().required("password  can't be empty."),
    });


    const handleSubmit = (values: any, helpers: any) => {

        console.log(values)

        AccountService.login(values.username, values.password).then(
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
                    helpers.setFieldError("password", "Account with this password does not exist.")
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
                    // console.log(values, helpers);
                    // helpers.setFieldError("username", "test")
                    // // handleSubmit(e)
                }}
                initialValues={{
                    username: '',
                    password: '',
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
                                <Form.Group as={Col} my="4">
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
                            </Row>
                        </Row>
                        <Row>
                            <Col className="h-100 d-flex justify-content-center align-items-center">
                                <Stack gap={4} direction="horizontal">
                                    <Button type="submit">sign up</Button>
                                    <Link to={"/register"}>
                                        <Button type="submit">I don't have an account</Button>
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

// return (
//     <Container className="mt-4">
// <Col className="h-100 d-flex justify-content-center align-items-center">
//     <h2>Sign in</h2>
// </Col>
//         <Form noValidate validated={validated} onSubmit={handleSubmit}>
//             <Row className="mb-4">
//                 <Col>
//                     <Stack gap={4} direction="vertical">
//                         <Stack gap={4} direction="vertical">
//                             <Form.Group controlId="Username">
//                                 <Form.Control
//                                     required
//                                     type="text"
//                                     placeholder="username"
//                                     ref={usernameRef}
//                                 />
//                                 <Form.Control.Feedback type="invalid">
//                                     Please choose a username.
//                                 </Form.Control.Feedback>
//                             </Form.Group>
//                             <Form.Group controlId="Password">
//                                 <Form.Control
//                                     required
//                                     type="password"
//                                     placeholder="password"
//                                     ref={passwordRef}
//                                 />
//                                 <Form.Control.Feedback type="invalid">
//                                     Please choose a password.
//                                 </Form.Control.Feedback>
//                             </Form.Group>
//                         </Stack>
//                     </Stack>
//                 </Col>
//             </Row>
//             <Row>
//                 <Col className="h-100 d-flex justify-content-center align-items-center">
//                     <Stack gap={4} direction="horizontal">
//                         <Button type="submit">sign in</Button>
//                         <Link to={"/register"}>
//                             <Button type="submit">sign up</Button>
//                         </Link>
//                     </Stack>
//                 </Col>
//             </Row>
//         </Form>
//     </Container>
// )
// }
export default LoginForm